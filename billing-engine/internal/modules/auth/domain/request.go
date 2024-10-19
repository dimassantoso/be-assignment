package domain

import (
	"billing-engine/pkg/helper"
	shareddomain "billing-engine/pkg/shared/domain"
)

// RequestAuth model
type RequestAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Deserialize to db model
func (r *RequestAuth) Deserialize() (res shareddomain.Auth) {
	res.Email = r.Email
	res.Password = helper.GeneratePassword(r.Password)

	return
}

// RequestLogin model
type RequestLogin struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	KeepSignIn bool   `json:"keep_sign_in"`
}
