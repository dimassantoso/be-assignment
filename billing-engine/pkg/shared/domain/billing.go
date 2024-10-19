package domain

import (
	"time"
)

// Billing model
type Billing struct {
	ID              int            `gorm:"column:id;primaryKey"`
	LoanID          int            `gorm:"column:loan_id;not null"`
	Week            int            `gorm:"column:week;not null"`
	DueDate         time.Time      `gorm:"column:due_date;not null"`
	AmountDue       float64        `gorm:"column:amount_due;not null"`
	PaymentDate     *time.Time     `gorm:"column:payment_date"`
	PaymentMethodID *int           `gorm:"column:payment_method_id"`
	CreatedAt       time.Time      `gorm:"column:created_at;default:now()"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;default:now()"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID"`
	Loan            *Loan          `gorm:"foreignKey:LoanID;references:ID"`
}

// TableName return table name of Billing model
func (Billing) TableName() string {
	return "billings"
}
