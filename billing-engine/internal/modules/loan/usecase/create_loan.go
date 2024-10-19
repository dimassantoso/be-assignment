package usecase

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"fmt"
	"github.com/golangid/candi/candihelper"
	"golang.org/x/sync/errgroup"
	"math"
	"time"

	borrowerdomain "billing-engine/internal/modules/borrower/domain"
	"billing-engine/internal/modules/loan/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *loanUsecaseImpl) CreateLoan(ctx context.Context, req *domain.RequestLoan) (result domain.ResponseLoan, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "LoanUsecase:CreateLoan")
	defer trace.Finish()

	var (
		eg, egCtx       = errgroup.WithContext(ctx)
		borrower        shareddomain.Borrower
		durationLoan    shareddomain.Duration
		createdLoan     shareddomain.Loan
		billingSchedule []shareddomain.Billing
	)

	eg.Go(func() error {
		borrower, err = uc.repoSQL.BorrowerRepo().Find(egCtx, &borrowerdomain.FilterBorrower{
			ID: &req.BorrowerID,
		})
		return err
	})

	eg.Go(func() error {
		durationLoan, err = uc.repoSQL.LoanRepo().FindDuration(egCtx, &domain.FilterDuration{
			ID: &req.DurationID,
		})
		return err
	})

	eg.Go(func() error {
		_, err = uc.repoSQL.LoanRepo().Find(egCtx, &domain.FilterLoan{
			BorrowerID:           &req.BorrowerID,
			IsContainOutstanding: candihelper.ToBoolPtr(true),
		})
		if err == nil {
			err = fmt.Errorf("borrower %v have active loan", req.BorrowerID)
			return err
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return
	}

	interestAmount := req.PrincipleAmount * durationLoan.Interest / 100
	totalAmount := req.PrincipleAmount + interestAmount
	installmentAmount := math.Ceil(totalAmount / float64(durationLoan.Week))

	err = uc.repoSQL.WithTransaction(ctx, func(ctx context.Context) error {
		createdLoan = shareddomain.Loan{
			BorrowerID:      req.BorrowerID,
			DurationID:      req.DurationID,
			PrincipalAmount: req.PrincipleAmount,
			InterestAmount:  interestAmount,
			TotalAmount:     totalAmount,
			LastPaymentDate: nil,
			Outstanding:     totalAmount,
		}
		err = uc.repoSQL.LoanRepo().Save(ctx, &createdLoan)
		if err != nil {
			return err
		}
		now := time.Now()
		dueDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
		for i := 1; i <= durationLoan.Week; i++ {
			dueDate = dueDate.AddDate(0, 0, 7)
			billingSchedule = append(billingSchedule, shareddomain.Billing{
				LoanID:    createdLoan.ID,
				Week:      i,
				DueDate:   dueDate,
				AmountDue: installmentAmount,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
		return uc.repoSQL.BillingRepo().SaveMany(ctx, billingSchedule)
	})
	if err != nil {
		return
	}

	result.Serialize(borrower, durationLoan, createdLoan, billingSchedule)
	return
}
