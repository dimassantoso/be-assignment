package usecase

import (
	"context"

	"billing-engine/internal/modules/billing/domain"
	"billing-engine/pkg/shared/repository"
	"billing-engine/pkg/shared/usecase/common"

	"github.com/golangid/candi/codebase/factory/dependency"
)

// BillingUsecase abstraction
type BillingUsecase interface {
	BillingRepayment(ctx context.Context, data *domain.RequestBillingRepayment) (err error)
}

type billingUsecaseImpl struct {
	deps          dependency.Dependency
	sharedUsecase common.Usecase
	repoSQL       repository.RepoSQL
}

// NewBillingUsecase usecase impl constructor
func NewBillingUsecase(deps dependency.Dependency) (BillingUsecase, func(sharedUsecase common.Usecase)) {
	uc := &billingUsecaseImpl{
		deps:    deps,
		repoSQL: repository.GetSharedRepoSQL(),
	}
	return uc, func(sharedUsecase common.Usecase) {
		uc.sharedUsecase = sharedUsecase
	}
}
