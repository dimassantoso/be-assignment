package domain

import (
	shareddomain "billing-engine/pkg/shared/domain"
	"github.com/golangid/candi/candishared"
	"time"
)

// ResponseLoanList model
type ResponseLoanList struct {
	Meta candishared.Meta `json:"meta"`
	Data []ResponseLoan   `json:"data"`
}

// ResponseLoan model
type ResponseLoan struct {
	ID              int               `json:"id"`
	BorrowerID      int               `json:"borrower_id"`
	BorrowerName    string            `json:"borrower_name"`
	PrincipalAmount float64           `json:"principal_amount"`
	InterestRate    float64           `json:"interest_rate"`
	InterestAmount  float64           `json:"interest_amount"`
	TotalAmount     float64           `json:"total_amount"`
	Duration        int               `json:"duration"`
	CreatedAt       string            `json:"createdAt"`
	UpdatedAt       string            `json:"updatedAt"`
	BillingSchedule []BillingSchedule `json:"billing_schedule"`
}

type BillingSchedule struct {
	Week      int     `json:"week"`
	DueDate   string  `json:"due_date"`
	AmountDue float64 `json:"amount_due"`
}

func (r *ResponseLoan) Serialize(borrower shareddomain.Borrower, durationLoan shareddomain.Duration, loan shareddomain.Loan, Billing []shareddomain.Billing) {
	r.ID = loan.ID
	r.BorrowerID = borrower.ID
	r.BorrowerName = borrower.Name
	r.PrincipalAmount = loan.PrincipalAmount
	r.InterestAmount = loan.InterestAmount
	r.TotalAmount = loan.TotalAmount
	r.Duration = durationLoan.Week
	r.InterestRate = durationLoan.Interest
	r.CreatedAt = loan.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = loan.UpdatedAt.Format(time.RFC3339)
	for _, b := range Billing {
		r.BillingSchedule = append(r.BillingSchedule, BillingSchedule{
			Week:      b.Week,
			DueDate:   b.DueDate.Format(time.DateOnly),
			AmountDue: b.AmountDue,
		})
	}
}

type ResponseLoanSimulation struct {
	Week      int     `json:"week"`
	AmountDue float64 `json:"amount_due"`
}

type ResponseLoanOutstanding struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Outstanding float64 `json:"outstanding"`
}
