package domain

import (
	"billing-engine/pkg/shared"
)

// FilterLoan model
type FilterLoan struct {
	shared.Filter
	ID                   *int     `json:"id"`
	BorrowerID           *int     `json:"borrower_id"`
	IsContainOutstanding *bool    `json:"is_contain_outstanding"`
	Preloads             []string `json:"-"`
}

// FilterDuration model
type FilterDuration struct {
	ID *int `json:"id"`
}
