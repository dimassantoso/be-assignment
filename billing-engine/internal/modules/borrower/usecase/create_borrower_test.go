package usecase

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"gorm.io/gorm"

	"billing-engine/internal/modules/borrower/domain"
	mockrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBorrowerUsecaseImpl_CreateBorrower(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {
		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{}, gorm.ErrRecordNotFound)
		borrowerRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.CreateBorrower(context.Background(), &domain.RequestBorrower{})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative", func(t *testing.T) {
		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{ID: 3}, nil)
		borrowerRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.CreateBorrower(context.Background(), &domain.RequestBorrower{})
		assert.Error(t, err)
	})
}
