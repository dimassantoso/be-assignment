package usecase

import (
	"context"

	"billing-engine/internal/modules/borrower/domain"
	"billing-engine/pkg/shared/repository"
	"billing-engine/pkg/shared/usecase/common"

	"github.com/golangid/candi/codebase/factory/dependency"
)

// BorrowerUsecase abstraction
type BorrowerUsecase interface {
	GetAllBorrower(ctx context.Context, filter *domain.FilterBorrower) (data domain.ResponseBorrowerList, err error)
	GetDetailBorrower(ctx context.Context, id int) (data domain.ResponseBorrower, err error)
	CreateBorrower(ctx context.Context, data *domain.RequestBorrower) (res domain.ResponseBorrower, err error)
	UpdateBorrower(ctx context.Context, data *domain.RequestBorrower) (err error)
	DeleteBorrower(ctx context.Context, id int) (err error)
	DelinquentCheck(ctx context.Context, id int) (res domain.ResponseDelinquentCheck, err error)
}

type borrowerUsecaseImpl struct {
	deps          dependency.Dependency
	sharedUsecase common.Usecase
	repoSQL       repository.RepoSQL
}

// NewBorrowerUsecase usecase impl constructor
func NewBorrowerUsecase(deps dependency.Dependency) (BorrowerUsecase, func(sharedUsecase common.Usecase)) {
	uc := &borrowerUsecaseImpl{
		deps:    deps,
		repoSQL: repository.GetSharedRepoSQL(),
	}
	return uc, func(sharedUsecase common.Usecase) {
		uc.sharedUsecase = sharedUsecase
	}
}
