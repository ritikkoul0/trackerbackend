package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	UserID        primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	Month         string             `bson:"month" json:"month"` // Format: "2024-01"
	Income        float64            `bson:"income" json:"income"`
	TotalExpenses float64            `bson:"total_expenses" json:"total_expenses"`
	Savings       float64            `bson:"savings" json:"savings"`
	SavingsGoal   float64            `bson:"savings_goal" json:"savings_goal"`
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


