package usecase

import (
	"context"
	"github.com/golangid/candi/codebase/interfaces"

	"billing-engine/internal/modules/auth/domain"
	"billing-engine/pkg/shared/repository"
	"billing-engine/pkg/shared/usecase/common"

	"github.com/golangid/candi/codebase/factory/dependency"
)

// AuthUsecase abstraction
type AuthUsecase interface {
	LoginAuth(ctx context.Context, data *domain.RequestLogin) (res domain.ResponseLogin, attemptLogin int, err error)
	CreateAuth(ctx context.Context, data *domain.RequestAuth) (err error)
}

type authUsecaseImpl struct {
	deps          dependency.Dependency
	sharedUsecase common.Usecase
	repoSQL       repository.RepoSQL
	cache         interfaces.Cache
}

// NewAuthUsecase usecase impl constructor
func NewAuthUsecase(deps dependency.Dependency) (AuthUsecase, func(sharedUsecase common.Usecase)) {
	uc := &authUsecaseImpl{
		deps:    deps,
		repoSQL: repository.GetSharedRepoSQL(),
	}
	if redisPool := deps.GetRedisPool(); redisPool != nil {
		uc.cache = redisPool.Cache()
	}
	return uc, func(sharedUsecase common.Usecase) {
		uc.sharedUsecase = sharedUsecase
	}
}
