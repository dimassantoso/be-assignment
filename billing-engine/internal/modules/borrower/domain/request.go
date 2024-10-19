package domain

import shareddomain "billing-engine/pkg/shared/domain"

// RequestBorrower model
type RequestBorrower struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Deserialize to db model
func (r *RequestBorrower) Deserialize() (res shareddomain.Borrower) {
	res.ID = r.ID
	res.Email = r.Email
	res.Name = r.Name
	return
}
