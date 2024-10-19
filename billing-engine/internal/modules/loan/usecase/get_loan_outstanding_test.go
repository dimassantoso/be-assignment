package usecase

import (
	mockborrowerrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mockrepo "billing-engine/pkg/mocks/modules/loan/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestLoanUsecaseImpl_GetLoanOutstanding(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {
		borrowerRepo := &mockborrowerrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{ID: 10}, nil)

		loanRepo := &mockrepo.LoanRepository{}
		loanRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Loan{ID: 10, Outstanding: 1000000}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := loanUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetLoanOutstanding(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative", func(t *testing.T) {
		borrowerRepo := &mockborrowerrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{}, gorm.ErrRecordNotFound)

		loanRepo := &mockrepo.LoanRepository{}
		loanRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Loan{}, gorm.ErrRecordNotFound)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("LoanRepo").Return(loanRepo)
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := loanUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetLoanOutstanding(context.Background(), 1)
		assert.Error(t, err)
	})
}
