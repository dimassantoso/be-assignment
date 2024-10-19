package usecase

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"errors"
	"github.com/golangid/candi/candishared"
	"golang.org/x/sync/errgroup"
	"time"

	"billing-engine/internal/modules/billing/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *billingUsecaseImpl) BillingRepayment(ctx context.Context, req *domain.RequestBillingRepayment) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "BillingUsecase:BillingRepayment")
	defer trace.Finish()

	var (
		eg, egCtx     = errgroup.WithContext(ctx)
		billing       shareddomain.Billing
		paymentMethod shareddomain.PaymentMethod
	)
	eg.Go(func() error {
		billing, err = uc.repoSQL.BillingRepo().Find(egCtx, &domain.FilterBilling{ID: &req.BillingID, Preloads: []string{"Loan"}})
		return err
	})

	eg.Go(func() error {
		paymentMethod, err = uc.repoSQL.BillingRepo().FindPaymentMethod(egCtx, &domain.FilterPaymentMethod{ID: &req.PaymentMethodID})
		return err
	})

	if err = eg.Wait(); err != nil {
		return err
	}

	if billing.PaymentDate != nil {
		err = errors.New("billing already paid")
		return
	}

	loan := billing.Loan
	return uc.repoSQL.WithTransaction(ctx, func(ctx context.Context) error {
		today := time.Now()
		loan.LastPaymentDate = &today
		loan.Outstanding = loan.Outstanding - billing.AmountDue

		err = uc.repoSQL.LoanRepo().Save(ctx, loan, candishared.DBUpdateSetUpdatedFields("LastPaymentDate", "Outstanding"))
		if err != nil {
			return err
		}

		billing.PaymentMethodID = &paymentMethod.ID
		billing.PaymentDate = &today
		return uc.repoSQL.BillingRepo().Save(ctx, &billing, candishared.DBUpdateSetUpdatedFields("PaymentMethodID", "PaymentDate"))
	})
}
