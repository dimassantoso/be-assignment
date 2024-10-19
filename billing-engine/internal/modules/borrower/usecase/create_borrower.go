package usecase

import (
	"context"
	"errors"

	"billing-engine/internal/modules/borrower/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *borrowerUsecaseImpl) CreateBorrower(ctx context.Context, req *domain.RequestBorrower) (result domain.ResponseBorrower, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerUsecase:CreateBorrower")
	defer trace.Finish()

	data := req.Deserialize()
	if _, err = uc.repoSQL.BorrowerRepo().Find(ctx, &domain.FilterBorrower{Email: data.Email}); err == nil {
		return domain.ResponseBorrower{}, errors.New("email exists")
	}
	err = uc.repoSQL.BorrowerRepo().Save(ctx, &data)
	result.Serialize(&data)
	return
}
