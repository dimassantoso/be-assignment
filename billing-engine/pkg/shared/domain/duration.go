package domain

// Duration model
type Duration struct {
	ID       int     `gorm:"column:id;primaryKey"`
	Week     int     `gorm:"column:week;not null"`
	Interest float64 `gorm:"column:interest;not null"`
}

// TableName return table name of Billing model
func (Duration) TableName() string {
	return "durations"
}
