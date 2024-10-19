package repository

import (
	"context"
	"database/sql"
	"fmt"

	// @candi:repositoryImport
	authrepo "billing-engine/internal/modules/auth/repository"
	billingrepo "billing-engine/internal/modules/billing/repository"
	borrowerrepo "billing-engine/internal/modules/borrower/repository"
	loanrepo "billing-engine/internal/modules/loan/repository"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"billing-engine/pkg/shared"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	// RepoSQL abstraction
	RepoSQL interface {
		WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) (err error)

		// @candi:repositoryMethod
		AuthRepo() authrepo.AuthRepository
		BillingRepo() billingrepo.BillingRepository
		LoanRepo() loanrepo.LoanRepository
		BorrowerRepo() borrowerrepo.BorrowerRepository
	}

	repoSQLImpl struct {
		readDB, writeDB *gorm.DB

		// register all repository from modules
		// @candi:repositoryField
		authRepo     authrepo.AuthRepository
		billingRepo  billingrepo.BillingRepository
		loanRepo     loanrepo.LoanRepository
		borrowerRepo borrowerrepo.BorrowerRepository
	}
)

var (
	globalRepoSQL RepoSQL
)

// setSharedRepoSQL set the global singleton "RepoSQL" implementation
func setSharedRepoSQL(readDB, writeDB *sql.DB) {
	gormRead, err := gorm.Open(postgres.New(postgres.Config{
		Conn: readDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	gormWrite, err := gorm.Open(postgres.New(postgres.Config{
		Conn: writeDB,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	shared.AddGormCallbacks(gormRead)
	shared.AddGormCallbacks(gormWrite)

	globalRepoSQL = NewRepositorySQL(gormRead, gormWrite)
}

// GetSharedRepoSQL returns the global singleton "RepoSQL" implementation
func GetSharedRepoSQL() RepoSQL {
	return globalRepoSQL
}

// NewRepositorySQL constructor
func NewRepositorySQL(readDB, writeDB *gorm.DB) RepoSQL {

	return &repoSQLImpl{
		readDB: readDB, writeDB: writeDB,

		// @candi:repositoryConstructor
		authRepo:     authrepo.NewAuthRepoSQL(readDB, writeDB),
		billingRepo:  billingrepo.NewBillingRepoSQL(readDB, writeDB),
		loanRepo:     loanrepo.NewLoanRepoSQL(readDB, writeDB),
		borrowerRepo: borrowerrepo.NewBorrowerRepoSQL(readDB, writeDB),
	}
}

// WithTransaction run transaction for each repository with context, include handle canceled or timeout context
func (r *repoSQLImpl) WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "RepoSQL:Transaction")
	defer trace.Finish()

	tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB)
	if !ok {
		tx = r.writeDB.Begin()
		if tx.Error != nil {
			return tx.Error
		}

		defer func() {
			if err != nil {
				tx.Rollback()
				trace.SetError(err)
			} else {
				tx.Commit()
			}
		}()
		ctx = candishared.SetToContext(ctx, candishared.ContextKeySQLTransaction, tx)
	}

	errChan := make(chan error)
	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic: %v", r)
			}
			close(errChan)
		}()

		if err := txFunc(ctx); err != nil {
			errChan <- err
		}
	}(ctx)

	select {
	case <-ctx.Done():
		return fmt.Errorf("Canceled or timeout: %v", ctx.Err())
	case e := <-errChan:
		return e
	}
}

// @candi:repositoryImplementation
func (r *repoSQLImpl) AuthRepo() authrepo.AuthRepository {
	return r.authRepo
}

func (r *repoSQLImpl) BillingRepo() billingrepo.BillingRepository {
	return r.billingRepo
}

func (r *repoSQLImpl) LoanRepo() loanrepo.LoanRepository {
	return r.loanRepo
}

func (r *repoSQLImpl) BorrowerRepo() borrowerrepo.BorrowerRepository {
	return r.borrowerRepo
}
