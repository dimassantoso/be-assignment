package domain

import (
	"billing-engine/pkg/shared"
	shareddomain "billing-engine/pkg/shared/domain"
	"github.com/golangid/candi/candihelper"
	"time"
)

// ResponseBorrowerList model
type ResponseBorrowerList struct {
	Meta shared.Meta        `json:"meta"`
	Data []ResponseBorrower `json:"data"`
}

// ResponseBorrower model
type ResponseBorrower struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	IsDelinquent bool          `json:"is_delinquent"`
	ActiveLoan   *BorrowerLoan `json:"active_loan,omitempty"`
}

type BorrowerLoan struct {
	ID              int           `json:"id"`
	BorrowerID      int           `json:"borrower_id"`
	PrincipalAmount float64       `json:"principal_amount"`
	InterestRate    float64       `json:"interest_rate"`
	InterestAmount  float64       `json:"interest_amount"`
	TotalAmount     float64       `json:"total_amount"`
	Duration        int           `json:"duration"`
	LastPaymentDate *string       `json:"last_payment_date"`
	Outstanding     float64       `json:"outstanding"`
	CreatedAt       string        `json:"created_at"`
	UpdatedAt       string        `json:"updated_at"`
	BillingPlan     []BillingPlan `json:"billing_plan"`
}

type BillingPlan struct {
	ID                int     `json:"id"`
	Week              int     `json:"week"`
	AmountDue         float64 `json:"amount_due"`
	DueDate           string  `json:"due_date"`
	PaymentDate       *string `json:"payment_date"`
	PaymentMethodID   *int    `json:"payment_method_id"`
	PaymentMethodName *string `json:"payment_method_name"`
	Status            string  `json:"status"`
}

// Serialize from db model
func (r *ResponseBorrower) Serialize(source *shareddomain.Borrower) {
	r.ID = source.ID
	r.Name = source.Name
	r.Email = source.Email
	r.CreatedAt = source.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = source.UpdatedAt.Format(time.RFC3339)

	countOverdueBilling := 0
	if source.ActiveLoan != nil {
		var lastPaymentData *string
		if source.ActiveLoan.LastPaymentDate != nil {
			lastPaymentData = candihelper.ToStringPtr(source.ActiveLoan.LastPaymentDate.Format(time.RFC3339))
		}
		borrowerLoan := BorrowerLoan{
			ID:              source.ActiveLoan.ID,
			BorrowerID:      source.ActiveLoan.BorrowerID,
			PrincipalAmount: source.ActiveLoan.PrincipalAmount,
			InterestRate:    source.ActiveLoan.Duration.Interest,
			InterestAmount:  source.ActiveLoan.InterestAmount,
			TotalAmount:     source.ActiveLoan.TotalAmount,
			Duration:        source.ActiveLoan.Duration.Week,
			LastPaymentDate: lastPaymentData,
			Outstanding:     source.ActiveLoan.Outstanding,
			CreatedAt:       source.ActiveLoan.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       source.ActiveLoan.UpdatedAt.Format(time.RFC3339),
		}

		for _, b := range source.ActiveLoan.Billing {
			var (
				paymentMethodID                *int
				paymentDate, paymentMethodName *string
				billingStatus                  string
			)
			if b.PaymentDate != nil {
				billingStatus = "PAID"
				paymentDate = candihelper.ToStringPtr(b.PaymentDate.Format(time.RFC3339))
				paymentMethodID = b.PaymentMethodID
				paymentMethodName = &b.PaymentMethod.Name
			} else {
				if b.DueDate.Before(time.Now()) {
					billingStatus = "OVERDUE"
					countOverdueBilling++
				} else {
					billingStatus = "UNPAID"
				}
			}
			borrowerLoan.BillingPlan = append(borrowerLoan.BillingPlan, BillingPlan{
				ID:                b.ID,
				Week:              b.Week,
				AmountDue:         b.AmountDue,
				DueDate:           b.DueDate.Format(time.RFC3339),
				PaymentDate:       paymentDate,
				PaymentMethodID:   paymentMethodID,
				PaymentMethodName: paymentMethodName,
				Status:            billingStatus,
			})
		}

		r.ActiveLoan = &borrowerLoan
		r.IsDelinquent = countOverdueBilling > 2
	}
}

type ResponseDelinquentCheck struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IsDelinquent bool   `json:"is_delinquent"`
}
