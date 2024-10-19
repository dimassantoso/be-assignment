package repository

import (
	"context"

	"billing-engine/internal/modules/auth/domain"
	shareddomain "billing-engine/pkg/shared/domain"
	"time"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"billing-engine/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type authRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewAuthRepoSQL mongo repo constructor
func NewAuthRepoSQL(readDB, writeDB *gorm.DB) AuthRepository {
	return &authRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *authRepoSQL) Find(ctx context.Context, filter *domain.FilterAuth) (result shareddomain.Auth, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "AuthRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterAuth(shared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *authRepoSQL) Save(ctx context.Context, data *shareddomain.Auth, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "AuthRepoSQL:Save")
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

func (r *authRepoSQL) setFilterAuth(db *gorm.DB, filter *domain.FilterAuth) *gorm.DB {

	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}

	if filter.Email != "" {
		db = db.Where("email = ?", filter.Email)
	}

	return db
}
