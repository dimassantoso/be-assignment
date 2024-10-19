package repository

import (
	"billing-engine/internal/modules/loan/domain"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"

	"github.com/golangid/candi/candishared"
)

// LoanRepository abstract interface
type LoanRepository interface {
	Find(ctx context.Context, filter *domain.FilterLoan) (shareddomain.Loan, error)
	Save(ctx context.Context, data *shareddomain.Loan, updateOptions ...candishared.DBUpdateOptionFunc) error
	FindDuration(ctx context.Context, filter *domain.FilterDuration) (shareddomain.Duration, error)
}
