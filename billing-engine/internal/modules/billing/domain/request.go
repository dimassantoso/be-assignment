package domain

// RequestBillingRepayment model
type RequestBillingRepayment struct {
	BillingID       int `json:"billing_id"`
	PaymentMethodID int `json:"payment_method_id"`
}
