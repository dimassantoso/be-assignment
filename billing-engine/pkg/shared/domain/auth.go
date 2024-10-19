package domain

import (
	"time"
)

// Auth model
type Auth struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Email     string    `gorm:"column:email;type:varchar(50)" json:"email"`
	Password  string    `gorm:"column:password;type:varchar(255)" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName return table name of Auth model
func (Auth) TableName() string {
	return "auths"
}
