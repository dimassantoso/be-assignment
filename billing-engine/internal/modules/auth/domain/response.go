package domain

// ResponseLogin model
type ResponseLogin struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
