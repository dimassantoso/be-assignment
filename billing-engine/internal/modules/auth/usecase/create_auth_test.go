package usecase

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"gorm.io/gorm"
	"time"

	"billing-engine/internal/modules/auth/domain"
	mockrepo "billing-engine/pkg/mocks/modules/auth/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUsecaseImpl_CreateAuth(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{}, gorm.ErrRecordNotFound)
		authRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.CreateAuth(context.Background(), &domain.RequestAuth{
			Email:    "admin@example.com",
			Password: "abc123",
		})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Negative (email already exists)", func(t *testing.T) {

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{
			ID:        1,
			Email:     "admin@example.com",
			Password:  "<hashed>",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.CreateAuth(context.Background(), &domain.RequestAuth{
			Email:    "admin@example.com",
			Password: "abc123",
		})
		assert.Error(t, err)
	})

	t.Run("Testcase #3: Negative (failed to create data to db)", func(t *testing.T) {

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{}, gorm.ErrRecordNotFound)
		authRepo.On("Save", mock.Anything, mock.Anything).Return(gorm.ErrInvalidData)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
		}

		err := uc.CreateAuth(context.Background(), &domain.RequestAuth{
			Email:    "admin@example.com",
			Password: "abc123",
		})
		assert.Error(t, err)
	})
}
