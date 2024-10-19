package usecase

import (
	borrowerdomain "billing-engine/internal/modules/borrower/domain"
	"billing-engine/internal/modules/loan/domain"
	"billing-engine/pkg/shared"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/tracer"
	"golang.org/x/sync/errgroup"
)

func (uc *loanUsecaseImpl) GetLoanOutstanding(ctx context.Context, borrowerID int) (result domain.ResponseLoanOutstanding, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "LoanUsecase:GetLoanSimulation")
	defer trace.Finish()

	var (
		eg, egCtx = errgroup.WithContext(ctx)
		borrower  shareddomain.Borrower
		loan      shareddomain.Loan
	)

	eg.Go(func() error {
		borrower, err = uc.repoSQL.BorrowerRepo().Find(egCtx, &borrowerdomain.FilterBorrower{
			ID: &borrowerID,
		})
		return err
	})

	eg.Go(func() error {
		loan, err = uc.repoSQL.LoanRepo().Find(egCtx, &domain.FilterLoan{
			Filter:               shared.Filter{ShowAll: true},
			BorrowerID:           &borrowerID,
			IsContainOutstanding: candihelper.ToBoolPtr(true),
		})
		return err
	})

	if err = eg.Wait(); err != nil {
		return
	}

	result = domain.ResponseLoanOutstanding{
		ID:          borrower.ID,
		Name:        borrower.Name,
		Outstanding: loan.Outstanding,
	}

	return
}
