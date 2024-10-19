package repository

import (
	"context"

	"time"

	"billing-engine/internal/modules/loan/domain"
	shareddomain "billing-engine/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"billing-engine/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type loanRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewLoanRepoSQL mongo repo constructor
func NewLoanRepoSQL(readDB, writeDB *gorm.DB) LoanRepository {
	return &loanRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *loanRepoSQL) Find(ctx context.Context, filter *domain.FilterLoan) (result shareddomain.Loan, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "LoanRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterLoan(shared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *loanRepoSQL) Save(ctx context.Context, data *shareddomain.Loan, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "LoanRepoSQL:Save")
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

func (r *loanRepoSQL) setFilterLoan(db *gorm.DB, filter *domain.FilterLoan) *gorm.DB {

	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}

	if filter.BorrowerID != nil {
		db = db.Where("borrower_id = ?", *filter.BorrowerID)
	}

	if filter.IsContainOutstanding != nil {
		isContainOutstanding := *filter.IsContainOutstanding
		if isContainOutstanding {
			db = db.Where("outstanding > 0")
		} else {
			db = db.Where("outstanding = 0")
		}
	}

	for _, preload := range filter.Preloads {
		db = db.Preload(preload)
	}

	return db
}

func (r *loanRepoSQL) FindDuration(ctx context.Context, filter *domain.FilterDuration) (result shareddomain.Duration, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "LoanRepoSQL:FindDuration")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterDuration(shared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *loanRepoSQL) setFilterDuration(db *gorm.DB, filter *domain.FilterDuration) *gorm.DB {

	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}

	return db
}
