package usecase

import (
	"context"
	
	"billing-engine/internal/modules/borrower/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *borrowerUsecaseImpl) DeleteBorrower(ctx context.Context, id int) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerUsecase:DeleteBorrower")
	defer trace.Finish()

	repoFilter := domain.FilterBorrower{ID: &id}
	return uc.repoSQL.BorrowerRepo().Delete(ctx, &repoFilter)
}
