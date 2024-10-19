package usecase

import (
	"context"

	"billing-engine/internal/modules/loan/domain"
	"billing-engine/pkg/shared/repository"
	"billing-engine/pkg/shared/usecase/common"

	"github.com/golangid/candi/codebase/factory/dependency"
)

// LoanUsecase abstraction
type LoanUsecase interface {
	CreateLoan(ctx context.Context, data *domain.RequestLoan) (res domain.ResponseLoan, err error)
	GetLoanSimulation(ctx context.Context, data *domain.RequestLoanSimulation) (res []domain.ResponseLoanSimulation, err error)
	GetLoanOutstanding(ctx context.Context, borrowerID int) (res domain.ResponseLoanOutstanding, err error)
}

type loanUsecaseImpl struct {
	deps          dependency.Dependency
	sharedUsecase common.Usecase
	repoSQL       repository.RepoSQL
}

// NewLoanUsecase usecase impl constructor
func NewLoanUsecase(deps dependency.Dependency) (LoanUsecase, func(sharedUsecase common.Usecase)) {
	uc := &loanUsecaseImpl{
		deps:    deps,
		repoSQL: repository.GetSharedRepoSQL(),
	}
	return uc, func(sharedUsecase common.Usecase) {
		uc.sharedUsecase = sharedUsecase
	}
}
