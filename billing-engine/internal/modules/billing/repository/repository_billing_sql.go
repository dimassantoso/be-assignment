package repository

import (
	"context"

	"time"

	"billing-engine/internal/modules/billing/domain"
	shareddomain "billing-engine/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"billing-engine/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type billingRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewBillingRepoSQL mongo repo constructor
func NewBillingRepoSQL(readDB, writeDB *gorm.DB) BillingRepository {
	return &billingRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *billingRepoSQL) Find(ctx context.Context, filter *domain.FilterBilling) (result shareddomain.Billing, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BillingRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterBilling(shared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *billingRepoSQL) Save(ctx context.Context, data *shareddomain.Billing, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BillingRepoSQL:Save")
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

func (r *billingRepoSQL) setFilterBilling(db *gorm.DB, filter *domain.FilterBilling) *gorm.DB {

	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}

	for _, preload := range filter.Preloads {
		db = db.Preload(preload)
	}

	return db
}

func (r *billingRepoSQL) SaveMany(ctx context.Context, data []shareddomain.Billing) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BillingRepoSQL:SaveMany")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.writeDB
	if tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB); ok {
		db = tx
	}

	return shared.SetSpanToGorm(ctx, db).Omit(clause.Associations).Create(data).Error
}

func (r *billingRepoSQL) FindPaymentMethod(ctx context.Context, filter *domain.FilterPaymentMethod) (result shareddomain.PaymentMethod, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BillingRepoSQL:FindPaymentMethod")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterPaymentMethod(shared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *billingRepoSQL) setFilterPaymentMethod(db *gorm.DB, filter *domain.FilterPaymentMethod) *gorm.DB {

	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}

	return db
}

func (r *billingRepoSQL) FindOverdueBilling(ctx context.Context, filter *domain.FilterOverdueBilling) (result shareddomain.OverdueBilling, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BillingRepoSQL:FindOverdueBilling")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.readDB
	db = db.Select(`b.id AS borrower_id, b.name AS borrower_name, COUNT(billings.id) AS count_overdue`)
	db = db.Joins(`JOIN loans l ON billings.loan_id = l.id AND l.outstanding > 0`)
	db = db.Joins(`JOIN borrowers b ON l.borrower_id = b.id`)
	db = db.Table("billings")
	db = db.Where("billings.payment_date IS NULL AND billings.due_date < NOW()")
	db = db.Group("b.id, b.name")

	err = r.setFilterOverdueBilling(shared.SetSpanToGorm(ctx, db), filter).Take(&result).Error
	return
}

func (r *billingRepoSQL) setFilterOverdueBilling(db *gorm.DB, filter *domain.FilterOverdueBilling) *gorm.DB {

	if filter.BorrowerID != nil {
		db = db.Where("borrower_id = ?", *filter.BorrowerID)
	}

	return db
}
