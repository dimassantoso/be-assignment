package usecase

import (
	"context"
	"errors"

	"billing-engine/internal/modules/borrower/domain"
	mockrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBorrowerUsecaseImpl_GetAllBorrower(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("FetchAll", mock.Anything, mock.Anything, mock.Anything).Return([]shareddomain.Borrower{
			{
				ID: 1,
			},
		}, nil)
		borrowerRepo.On("Count", mock.Anything, mock.Anything).Return(10)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetAllBorrower(context.Background(), &domain.FilterBorrower{})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("FetchAll", mock.Anything, mock.Anything, mock.Anything).Return([]shareddomain.Borrower{}, errors.New("Error"))
		borrowerRepo.On("Count", mock.Anything, mock.Anything).Return(10)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetAllBorrower(context.Background(), &domain.FilterBorrower{})
		assert.Error(t, err)
	})
}
