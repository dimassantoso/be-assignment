package usecase

import (
	"context"
	"fmt"

	"billing-engine/internal/modules/auth/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *authUsecaseImpl) CreateAuth(ctx context.Context, req *domain.RequestAuth) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "AuthUsecase:CreateAuth")
	defer trace.Finish()

	_, err := uc.repoSQL.AuthRepo().Find(ctx, &domain.FilterAuth{
		Email: req.Email,
	})
	if err == nil {
		err = fmt.Errorf("email %s already exists", req.Email)
		return err
	}
	data := req.Deserialize()
	return uc.repoSQL.AuthRepo().Save(ctx, &data)
}
