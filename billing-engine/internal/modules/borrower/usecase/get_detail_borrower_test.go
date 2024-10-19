package usecase

import (
	"context"
	"gorm.io/gorm"

	mockrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBorrowerUsecaseImpl_GetDetailBorrower(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetDetailBorrower(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{}, gorm.ErrRecordNotFound)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.GetDetailBorrower(context.Background(), 1)
		assert.Error(t, err)
	})
}
