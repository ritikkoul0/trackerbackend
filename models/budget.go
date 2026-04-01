package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID        uint    `gorm:"not null;index" json:"user_id"`
	User          User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Month         string  `gorm:"type:varchar(7);not null;index" json:"month"` // Format: "2024-01"
	Income        float64 `gorm:"type:decimal(15,2);default:0" json:"income"`
	TotalExpenses float64 `gorm:"type:decimal(15,2);default:0" json:"total_expenses"`
	Savings       float64 `gorm:"type:decimal(15,2);default:0" json:"savings"`
	SavingsGoal   float64 `gorm:"type:decimal(15,2);default:0" json:"savings_goal"`
}

// CalculateSavings calculates savings from income and expenses
func (b *Budget) CalculateSavings() {
	b.Savings = b.Income - b.TotalExpenses
}

// CalculateSavingsPercentage calculates the savings percentage
func (b *Budget) CalculateSavingsPercentage() float64 {
	if b.SavingsGoal > 0 {
		return (b.Savings / b.SavingsGoal) * 100
	}
	return 0
}

//
