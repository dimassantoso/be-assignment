package usecase

import (
	"billing-engine/pkg/shared"
	"context"

	"billing-engine/internal/modules/borrower/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *borrowerUsecaseImpl) GetAllBorrower(ctx context.Context, filter *domain.FilterBorrower) (result domain.ResponseBorrowerList, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerUsecase:GetAllBorrower")
	defer trace.Finish()

	data, err := uc.repoSQL.BorrowerRepo().FetchAll(ctx, filter)
	if err != nil {
		return result, err
	}
	count := uc.repoSQL.BorrowerRepo().Count(ctx, filter)
	result.Meta = shared.NewMeta(filter.Page, filter.Limit, count)

	result.Data = make([]domain.ResponseBorrower, len(data))
	for i, detail := range data {
		result.Data[i].Serialize(&detail)
	}

	return
}
