package repository

import (
	"context"

	"billing-engine/internal/modules/billing/domain"
	shareddomain "billing-engine/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// BillingRepository abstract interface
type BillingRepository interface {
	Find(ctx context.Context, filter *domain.FilterBilling) (shareddomain.Billing, error)
	Save(ctx context.Context, data *shareddomain.Billing, updateOptions ...candishared.DBUpdateOptionFunc) error
	SaveMany(ctx context.Context, data []shareddomain.Billing) error
	FindPaymentMethod(ctx context.Context, filter *domain.FilterPaymentMethod) (shareddomain.PaymentMethod, error)
	FindOverdueBilling(ctx context.Context, filter *domain.FilterOverdueBilling) (shareddomain.OverdueBilling, error)
}
