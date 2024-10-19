package domain

import (
	"billing-engine/pkg/shared"
)

// FilterBilling model
type FilterBilling struct {
	shared.Filter
	ID       *int     `json:"id"`
	Preloads []string `json:"-"`
}

// FilterPaymentMethod model
type FilterPaymentMethod struct {
	ID *int `json:"id"`
}

type FilterOverdueBilling struct {
	BorrowerID *int `json:"borrower_id"`
}
