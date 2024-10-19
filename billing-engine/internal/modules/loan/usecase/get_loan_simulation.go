package usecase

import (
	"billing-engine/internal/modules/loan/domain"
	"context"
	"github.com/golangid/candi/tracer"
	"math"
)

func (uc *loanUsecaseImpl) GetLoanSimulation(ctx context.Context, req *domain.RequestLoanSimulation) (result []domain.ResponseLoanSimulation, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "LoanUsecase:GetLoanSimulation")
	defer trace.Finish()

	duration, err := uc.repoSQL.LoanRepo().FindDuration(ctx, &domain.FilterDuration{ID: &req.DurationID})
	if err != nil {
		return
	}

	interestAmount := req.PrincipleAmount * duration.Interest / 100
	totalAmount := req.PrincipleAmount + interestAmount
	installmentAmount := math.Ceil(totalAmount / float64(duration.Week))

	for i := 1; i <= duration.Week; i++ {
		result = append(result, domain.ResponseLoanSimulation{
			Week:      i,
			AmountDue: installmentAmount,
		})
	}

	return
}
