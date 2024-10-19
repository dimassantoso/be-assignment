package usecase

import (
	"billing-engine/internal/modules/borrower/domain"
	mockrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"fmt"
	"github.com/golangid/candi/candihelper"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBorrowerUsecaseImpl_UpdateBorrower(t *testing.T) {
	ctx := context.Background()
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{ID: 1}, nil)
		borrowerRepo.On("Save", mock.Anything, mock.Anything, mock.AnythingOfType("candishared.DBUpdateOptionFunc")).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)
		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.UpdateBorrower(ctx, &domain.RequestBorrower{ID: 1})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Email already exists", func(t *testing.T) {
		borrowerRepo := new(mockrepo.BorrowerRepository)
		borrowerRepo.On("Find", mock.Anything, &domain.FilterBorrower{ID: candihelper.ToIntPtr(1)}).Return(shareddomain.Borrower{ID: 1}, nil).Once()
		borrowerRepo.On("Find", mock.Anything, &domain.FilterBorrower{Email: "existing@example.com"}).Return(shareddomain.Borrower{ID: 2}, nil).Once()

		repoSQL := new(mocksharedrepo.RepoSQL)
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		data := &domain.RequestBorrower{
			ID:    1,
			Email: "existing@example.com",
		}

		err := uc.UpdateBorrower(ctx, data)
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email %s already exists", data.Email), err)
	})

	t.Run("Testcase #3: Other errors during Find", func(t *testing.T) {
		borrowerRepo := new(mockrepo.BorrowerRepository)
		borrowerRepo.On("Find", mock.Anything, &domain.FilterBorrower{ID: candihelper.ToIntPtr(1)}).Return(shareddomain.Borrower{}, gorm.ErrRecordNotFound).Once()
		borrowerRepo.On("Find", mock.Anything, &domain.FilterBorrower{Email: "existing@example.com"}).Return(shareddomain.Borrower{}, nil).Once()

		repoSQL := new(mocksharedrepo.RepoSQL)
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		data := &domain.RequestBorrower{
			ID:    1,
			Email: "existing@example.com",
		}

		err := uc.UpdateBorrower(ctx, data)

		assert.Error(t, err)
	})
}
