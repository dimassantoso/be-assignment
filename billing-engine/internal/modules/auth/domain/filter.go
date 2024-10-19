package domain

import (
	"billing-engine/pkg/shared"
)

// FilterAuth model
type FilterAuth struct {
	shared.Filter
	ID       *int     `json:"id"`
	Email    string   `json:"email"`
	Preloads []string `json:"-"`
}
