package domain

// PaymentMethod model
type PaymentMethod struct {
	ID   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name;not null"`
}

// TableName return table name of Billing model
func (PaymentMethod) TableName() string {
	return "payment_methods"
}
