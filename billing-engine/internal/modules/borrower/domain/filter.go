package domain

import (
	"billing-engine/pkg/shared"
)

// FilterBorrower model
type FilterBorrower struct {
	shared.Filter
	ID       *int     `json:"id"`
	Email    string   `json:"email"`
	Preloads []string `json:"-"`
}
