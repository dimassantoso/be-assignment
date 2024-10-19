package repository

import (
	"context"

	"strings"
	"time"

	"billing-engine/internal/modules/borrower/domain"
	shareddomain "billing-engine/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"billing-engine/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type borrowerRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewBorrowerRepoSQL mongo repo constructor
func NewBorrowerRepoSQL(readDB, writeDB *gorm.DB) BorrowerRepository {
	return &borrowerRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *borrowerRepoSQL) FetchAll(ctx context.Context, filter *domain.FilterBorrower) (data []shareddomain.Borrower, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerRepoSQL:FetchAll")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	if filter.OrderBy == "" {
		filter.OrderBy = "updated_at"
	}

	db := r.setFilterBorrower(shared.SetSpanToGorm(ctx, r.readDB), filter).Order(clause.OrderByColumn{
		Column: clause.Column{Name: filter.OrderBy},
		Desc:   strings.ToUpper(filter.Sort) == "DESC",
	})
	if filter.Limit > 0 || !filter.ShowAll {
		db = db.Limit(filter.Limit).Offset(filter.CalculateOffset())
	}
	err = db.Find(&data).Error
	return
}

func (r *borrowerRepoSQL) Count(ctx context.Context, filter *domain.FilterBorrower) (count int) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerRepoSQL:Count")
	defer trace.Finish()

	var total int64
	r.setFilterBorrower(shared.SetSpanToGorm(ctx, r.readDB), filter).Model(&shareddomain.Borrower{}).Count(&total)
	count = int(total)

	trace.Log("count", count)
	return
}

func (r *borrowerRepoSQL) Find(ctx context.Context, filter *domain.FilterBorrower) (result shareddomain.Borrower, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterBorrower(shared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *borrowerRepoSQL) Save(ctx context.Context, data *shareddomain.Borrower, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerRepoSQL:Save")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.writeDB
	if tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB); ok {
		db = tx
	}
	data.UpdatedAt = time.Now()
	if data.CreatedAt.IsZero() {
		data.CreatedAt = time.Now()
	}
	if data.ID == 0 {
		err = shared.SetSpanToGorm(ctx, db).Omit(clause.Associations).Create(data).Error
	} else {
		err = shared.SetSpanToGorm(ctx, db).Model(data).Omit(clause.Associations).Updates(r.updateTools.ToMap(data, updateOptions...)).Error
	}
	return
}

func (r *borrowerRepoSQL) Delete(ctx context.Context, filter *domain.FilterBorrower) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerRepoSQL:Delete")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.writeDB
	if tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB); ok {
		db = tx
	}
	err = r.setFilterBorrower(shared.SetSpanToGorm(ctx, db), filter).Model(&shareddomain.Borrower{}).Update("deleted_at", time.Now()).Error
	return
}

func (r *borrowerRepoSQL) setFilterBorrower(db *gorm.DB, filter *domain.FilterBorrower) *gorm.DB {

	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}
	if filter.Search != "" {
		db = db.Where("(email ILIKE ? OR name ILIKE ?)", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.Email != "" {
		db = db.Where("email = ?", filter.Email)
	}

	for _, preload := range filter.Preloads {
		switch preload {
		case "Loan":
			db = db.Preload(preload, func(db *gorm.DB) *gorm.DB {
				db.Order("id")
				return db
			})
		case "ActiveLoan":
			db = db.Preload(preload, func(db *gorm.DB) *gorm.DB {
				db.Order("id")
				return db
			})
		case "Loan.Billing", "ActiveLoan.Billing":
			db = db.Preload(preload, func(db *gorm.DB) *gorm.DB {
				db.Order("loan_id, week")
				return db
			})
		default:
			db = db.Preload(preload)
		}
	}

	return db
}
