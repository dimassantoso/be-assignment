package usecase

import (
	mockbillingrepo "billing-engine/pkg/mocks/modules/billing/repository"
	mockborrowerrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"gorm.io/gorm"

	"billing-engine/internal/modules/loan/domain"
	mockrepo "billing-engine/pkg/mocks/modules/loan/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_loanUsecaseImpl_CreateLoan(t *testing.T) {
	ctx := context.Background()
	t.Run("Testcase #1: Positive", func(t *testing.T) {
		borrowerRepo := &mockborrowerrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{ID: 1}, nil)

		loanRepo := &mockrepo.LoanRepository{}
		loanRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Loan{}, gorm.ErrRecordNotFound)
		loanRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
		loanRepo.On("FindDuration", mock.Anything, mock.Anything).Return(shareddomain.Duration{
			ID:       1,
			Week:     50,
			Interest: 10,
		}, nil)

		billingRepo := &mockbillingrepo.BillingRepository{}
		billingRepo.On("SaveMany", mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)
		repoSQL.On("BillingRepo").Return(billingRepo)
		repoSQL.On("WithTransaction", mock.Anything, mock.AnythingOfType(`func(context.Context) error`), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(func(context.Context) error)
			err := arg(ctx)
			if err != nil {
				return
			}
		})

		uc := loanUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.CreateLoan(context.Background(), &domain.RequestLoan{
			BorrowerID:      1,
			PrincipleAmount: 1000000000,
			DurationID:      10,
		})
		assert.NoError(t, err)
	})
}
