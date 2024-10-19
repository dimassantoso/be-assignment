package domain

// RequestLoan model
type RequestLoan struct {
	BorrowerID      int     `json:"borrower_id"`
	PrincipleAmount float64 `json:"principle_amount"`
	DurationID      int     `json:"duration_id"`
}

type RequestLoanSimulation struct {
	PrincipleAmount float64 `json:"principle_amount"`
	DurationID      int     `json:"duration_id"`
}
