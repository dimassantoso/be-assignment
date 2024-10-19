package usecase

import (
	billingdomain "billing-engine/internal/modules/billing/domain"
	"billing-engine/internal/modules/borrower/domain"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"errors"
	"github.com/golangid/candi/tracer"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func (uc *borrowerUsecaseImpl) DelinquentCheck(ctx context.Context, id int) (result domain.ResponseDelinquentCheck, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BorrowerUsecase:DelinquentCheck")
	defer trace.Finish()

	var (
		eg, egCtx      = errgroup.WithContext(ctx)
		borrower       shareddomain.Borrower
		overdueBilling shareddomain.OverdueBilling
	)

	eg.Go(func() error {
		borrower, err = uc.repoSQL.BorrowerRepo().Find(egCtx, &domain.FilterBorrower{ID: &id})
		return err
	})

	eg.Go(func() error {
		overdueBilling, err = uc.repoSQL.BillingRepo().FindOverdueBilling(egCtx, &billingdomain.FilterOverdueBilling{BorrowerID: &id})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})

	if err = eg.Wait(); err != nil {
		return
	}

	result = domain.ResponseDelinquentCheck{
		ID:           borrower.ID,
		Name:         borrower.Name,
		IsDelinquent: overdueBilling.CountOverdue > 2,
	}

	return
}
