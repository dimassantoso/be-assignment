package usecase

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"github.com/golangid/candi/candihelper"
	"gorm.io/gorm"
	"time"

	"billing-engine/internal/modules/billing/domain"
	mockrepo "billing-engine/pkg/mocks/modules/billing/repository"
	mockloanrepo "billing-engine/pkg/mocks/modules/loan/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBillingUsecaseImpl_BillingRepayment(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {
		ctx := context.Background()
		billingRepo := &mockrepo.BillingRepository{}
		billingRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Billing{
			ID:        1,
			LoanID:    1,
			Week:      1,
			DueDate:   time.Now(),
			AmountDue: 100000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Loan: &shareddomain.Loan{
				ID:          1,
				Outstanding: 5000000,
			},
		}, nil)
		billingRepo.On("FindPaymentMethod", mock.Anything, mock.Anything).Return(shareddomain.PaymentMethod{ID: 1, Name: "Cash"}, nil)
		billingRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		loanRepo := &mockloanrepo.LoanRepository{}
		loanRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BillingRepo").Return(billingRepo)
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("WithTransaction", mock.Anything, mock.AnythingOfType(`func(context.Context) error`), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(func(context.Context) error)
			err := arg(ctx)
			if err != nil {
				return
			}
		})

		uc := billingUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.BillingRepayment(context.Background(), &domain.RequestBillingRepayment{})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative", func(t *testing.T) {
		ctx := context.Background()
		billingRepo := &mockrepo.BillingRepository{}
		billingRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Billing{}, gorm.ErrRecordNotFound)
		billingRepo.On("FindPaymentMethod", mock.Anything, mock.Anything).Return(shareddomain.PaymentMethod{ID: 1, Name: "Cash"}, nil)
		billingRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		loanRepo := &mockloanrepo.LoanRepository{}
		loanRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BillingRepo").Return(billingRepo)
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("WithTransaction", mock.Anything, mock.AnythingOfType(`func(context.Context) error`), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(func(context.Context) error)
			err := arg(ctx)
			if err != nil {
				return
			}
		})

		uc := billingUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.BillingRepayment(context.Background(), &domain.RequestBillingRepayment{})
		assert.Error(t, err)
	})

	t.Run("Testcase #3: Negative", func(t *testing.T) {
		ctx := context.Background()
		billingRepo := &mockrepo.BillingRepository{}
		billingRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Billing{
			ID:        1,
			LoanID:    1,
			Week:      1,
			DueDate:   time.Now(),
			AmountDue: 100000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Loan: &shareddomain.Loan{
				ID:          1,
				Outstanding: 5000000,
			},
		}, nil)
		billingRepo.On("FindPaymentMethod", mock.Anything, mock.Anything).Return(shareddomain.PaymentMethod{ID: 1, Name: "Cash"}, nil)
		billingRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		loanRepo := &mockloanrepo.LoanRepository{}
		loanRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BillingRepo").Return(billingRepo)
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("WithTransaction", mock.Anything, mock.AnythingOfType(`func(context.Context) error`), mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Run(func(args mock.Arguments) {
			arg := args.Get(1).(func(context.Context) error)
			err := arg(ctx)
			if err != nil {
				return
			}
		})

		uc := billingUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.BillingRepayment(context.Background(), &domain.RequestBillingRepayment{})
		assert.Error(t, err)
	})

	t.Run("Testcase #4: Negative", func(t *testing.T) {
		ctx := context.Background()
		billingRepo := &mockrepo.BillingRepository{}
		billingRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Billing{
			ID:          1,
			LoanID:      1,
			Week:        1,
			DueDate:     time.Now(),
			AmountDue:   100000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PaymentDate: candihelper.ToTimePtr(time.Now()),
			Loan: &shareddomain.Loan{
				ID:          1,
				Outstanding: 5000000,
			},
		}, nil)
		billingRepo.On("FindPaymentMethod", mock.Anything, mock.Anything).Return(shareddomain.PaymentMethod{ID: 1, Name: "Cash"}, nil)
		billingRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		loanRepo := &mockloanrepo.LoanRepository{}
		loanRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BillingRepo").Return(billingRepo)
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("WithTransaction", mock.Anything, mock.AnythingOfType(`func(context.Context) error`), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(func(context.Context) error)
			err := arg(ctx)
			if err != nil {
				return
			}
		})

		uc := billingUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.BillingRepayment(context.Background(), &domain.RequestBillingRepayment{})
		assert.Error(t, err)
	})
}
