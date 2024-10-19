package domain

import (
	"time"
)

// Borrower model
type Borrower struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string     `gorm:"column:name;not null"`
	Email      string     `gorm:"column:email;unique;not null"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
	CreatedAt  time.Time  `gorm:"column:created_at;default:now()"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;default:now()"`
	ActiveLoan *Loan      `gorm:"foreignKey:BorrowerID;references:ID"`
	Loan       []Loan     `gorm:"foreignKey:BorrowerID;references:ID"`
}

// TableName return table name of Borrower model
func (Borrower) TableName() string {
	return "borrowers"
}

type OverdueBilling struct {
	ID           int    `gorm:"column:id"`
	Name         string `gorm:"column:name"`
	CountOverdue int    `gorm:"column:count_overdue"`
}
