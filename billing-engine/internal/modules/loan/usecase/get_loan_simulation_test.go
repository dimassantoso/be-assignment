package usecase

import (
	"billing-engine/internal/modules/loan/domain"
	mockrepo "billing-engine/pkg/mocks/modules/loan/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestLoanUsecaseImpl_GetLoanSimulation(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		loanRepo := &mockrepo.LoanRepository{}
		loanRepo.On("FindDuration", mock.Anything, mock.Anything).Return(shareddomain.Duration{
			ID:       1,
			Week:     50,
			Interest: 10,
		}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("LoanRepo").Return(loanRepo)

		uc := loanUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetLoanSimulation(context.Background(), &domain.RequestLoanSimulation{
			PrincipleAmount: 10000000,
			DurationID:      1,
		})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative", func(t *testing.T) {

		loanRepo := &mockrepo.LoanRepository{}
		loanRepo.On("FindDuration", mock.Anything, mock.Anything).Return(shareddomain.Duration{}, gorm.ErrRecordNotFound)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("LoanRepo").Return(loanRepo)

		uc := loanUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetLoanSimulation(context.Background(), &domain.RequestLoanSimulation{
			PrincipleAmount: 10000000,
			DurationID:      1,
		})
		assert.Error(t, err)
	})
}
