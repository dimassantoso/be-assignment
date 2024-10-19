package usecase

import (
	"context"

	mockrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBorrowerUsecaseImpl_DeleteBorrower(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.DeleteBorrower(context.Background(), 1)
		assert.NoError(t, err)
	})
}
