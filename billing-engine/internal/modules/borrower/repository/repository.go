package repository

import (
	"context"

	"billing-engine/internal/modules/borrower/domain"
	shareddomain "billing-engine/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// BorrowerRepository abstract interface
type BorrowerRepository interface {
	FetchAll(ctx context.Context, filter *domain.FilterBorrower) ([]shareddomain.Borrower, error)
	Count(ctx context.Context, filter *domain.FilterBorrower) int
	Find(ctx context.Context, filter *domain.FilterBorrower) (shareddomain.Borrower, error)
	Save(ctx context.Context, data *shareddomain.Borrower, updateOptions ...candishared.DBUpdateOptionFunc) error
	Delete(ctx context.Context, filter *domain.FilterBorrower) (err error)
}
