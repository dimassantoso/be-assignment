package usecase

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"

	"billing-engine/internal/modules/borrower/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *borrowerUsecaseImpl) UpdateBorrower(ctx context.Context, data *domain.RequestBorrower) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerUsecase:UpdateBorrower")
	defer trace.Finish()

	var (
		borrowerByID, borrowerByEmail shareddomain.Borrower
		eg, egCtx                     = errgroup.WithContext(ctx)
	)

	eg.Go(func() error {
		borrowerByID, err = uc.repoSQL.BorrowerRepo().Find(egCtx, &domain.FilterBorrower{ID: &data.ID})
		return err
	})

	eg.Go(func() error {
		borrowerByEmail, err = uc.repoSQL.BorrowerRepo().Find(egCtx, &domain.FilterBorrower{Email: data.Email})
		return err
	})

	if err = eg.Wait(); err != nil {
		return
	}

	if borrowerByID.ID != borrowerByEmail.ID {
		err = fmt.Errorf("email %s already exists", data.Email)
		return
	}

	updateData := data.Deserialize()
	return uc.repoSQL.BorrowerRepo().Save(ctx, &updateData, candishared.DBUpdateSetUpdatedFields("Name", "Email"))
}
