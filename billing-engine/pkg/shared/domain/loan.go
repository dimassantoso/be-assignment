package domain

import (
	"time"
)

// Loan model
type Loan struct {
	ID              int        `gorm:"column:id;primaryKey;autoIncrement"`
	BorrowerID      int        `gorm:"column:borrower_id;not null"`
	DurationID      int        `gorm:"column:duration_id;not null"`
	PrincipalAmount float64    `gorm:"column:principal_amount;not null"`
	InterestAmount  float64    `gorm:"column:interest_amount;not null"`
	TotalAmount     float64    `gorm:"column:total_amount;not null"`
	LastPaymentDate *time.Time `gorm:"column:last_payment_date"`
	Outstanding     float64    `gorm:"column:outstanding;not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;default:now()"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;default:now()"`
	Borrower        *Borrower  `gorm:"foreignKey:BorrowerID;references:ID"`
	Duration        *Duration  `gorm:"foreignKey:DurationID;references:ID"`
	Billing         []Billing  `gorm:"foreignKey:LoanID;references:ID"`
}

// TableName return table name of Loan model
func (Loan) TableName() string {
	return "loans"
}
