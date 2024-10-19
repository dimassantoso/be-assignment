package usecase

import (
	"billing-engine/internal/modules/auth/domain"
	"billing-engine/pkg/helper"
	mockrepo "billing-engine/pkg/mocks/modules/auth/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	codebasemocks "github.com/golangid/candi/mocks/codebase/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestAuthUsecaseImpl_LoginAuth(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {
		mockCache := new(codebasemocks.Cache)
		mockCache.On("Get", mock.Anything, mock.Anything).Return(nil, nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{
			ID:        1,
			Email:     "example@example.com",
			Password:  helper.GeneratePassword("12345678"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
			cache:   mockCache,
		}

		_, _, err := uc.LoginAuth(context.Background(), &domain.RequestLogin{
			Email:      "example@example.com",
			Password:   "12345678",
			KeepSignIn: true,
		})
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Positive", func(t *testing.T) {
		mockCache := new(codebasemocks.Cache)
		mockCache.On("Get", mock.Anything, mock.Anything).Return([]byte{3}, nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{
			ID:        1,
			Email:     "example@example.com",
			Password:  helper.GeneratePassword("12345678"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
			cache:   mockCache,
		}

		_, _, err := uc.LoginAuth(context.Background(), &domain.RequestLogin{
			Email:      "example@example.com",
			Password:   "12345678",
			KeepSignIn: true,
		})
		assert.NoError(t, err)
	})

	t.Run("Testcase #3: Negative", func(t *testing.T) {
		mockCache := new(codebasemocks.Cache)
		mockCache.On("Get", mock.Anything, mock.Anything).Return([]byte{3}, nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{
			ID:        1,
			Email:     "example@example.com",
			Password:  helper.GeneratePassword("12345678"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
			cache:   mockCache,
		}

		_, _, err := uc.LoginAuth(context.Background(), &domain.RequestLogin{
			Email:      "example@example.com",
			Password:   "1234567890",
			KeepSignIn: true,
		})
		assert.Error(t, err)
	})

	t.Run("Testcase #4: Negative", func(t *testing.T) {
		mockCache := new(codebasemocks.Cache)
		mockCache.On("Get", mock.Anything, mock.Anything).Return([]byte{3}, nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{}, gorm.ErrRecordNotFound)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
			cache:   mockCache,
		}

		_, _, err := uc.LoginAuth(context.Background(), &domain.RequestLogin{
			Email:      "example@example.com",
			Password:   "1234567890",
			KeepSignIn: true,
		})
		assert.Error(t, err)
	})

	t.Run("Testcase #5: Positive", func(t *testing.T) {
		mockCache := new(codebasemocks.Cache)
		mockCache.On("Get", mock.Anything, mock.Anything).Return([]byte{4}, nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		uc := authUsecaseImpl{
			cache: mockCache,
		}

		_, _, err := uc.LoginAuth(context.Background(), &domain.RequestLogin{
			Email:      "example@example.com",
			Password:   "12345678",
			KeepSignIn: true,
		})
		assert.Error(t, err)
	})

	t.Run("Testcase #6: Negative", func(t *testing.T) {
		mockCache := new(codebasemocks.Cache)
		mockCache.On("Get", mock.Anything, mock.Anything).Return([]byte{3}, nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		authRepo := &mockrepo.AuthRepository{}
		authRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Auth{}, gorm.ErrInvalidDB)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("AuthRepo").Return(authRepo)

		uc := authUsecaseImpl{
			repoSQL: repoSQL,
			cache:   mockCache,
		}

		_, _, err := uc.LoginAuth(context.Background(), &domain.RequestLogin{
			Email:      "example@example.com",
			Password:   "1234567890",
			KeepSignIn: true,
		})
		assert.Error(t, err)
	})
}
