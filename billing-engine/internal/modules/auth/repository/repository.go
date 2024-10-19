package repository

import (
	"context"

	"billing-engine/internal/modules/auth/domain"
	shareddomain "billing-engine/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// AuthRepository abstract interface
type AuthRepository interface {
	Find(ctx context.Context, filter *domain.FilterAuth) (shareddomain.Auth, error)
	Save(ctx context.Context, data *shareddomain.Auth, updateOptions ...candishared.DBUpdateOptionFunc) error
}
