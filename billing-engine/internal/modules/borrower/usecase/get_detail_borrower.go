package usecase

import (
	"context"

	"billing-engine/internal/modules/borrower/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *borrowerUsecaseImpl) GetDetailBorrower(ctx context.Context, id int) (result domain.ResponseBorrower, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerUsecase:GetDetailBorrower")
	defer trace.Finish()

	repoFilter := domain.FilterBorrower{ID: &id, Preloads: []string{"ActiveLoan", "ActiveLoan.Billing", "ActiveLoan.Duration", "ActiveLoan.Billing.PaymentMethod"}}
	data, err := uc.repoSQL.BorrowerRepo().Find(ctx, &repoFilter)
	if err != nil {
		return result, err
	}

	result.Serialize(&data)

	return
}
