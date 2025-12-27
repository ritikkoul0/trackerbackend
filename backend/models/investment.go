package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Investment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	UserID       primitive.ObjectID  `bson:"user_id,omitempty" json:"user_id,omitempty"`
	GoalID       *primitive.ObjectID `bson:"goal_id,omitempty" json:"goal_id,omitempty"` // Link to Goal
	Name         string              `bson:"name" json:"name"`
	Type         string              `bson:"type" json:"type"` // Stocks, Mutual Fund, ETF, FD, PPF, etc.
	Invested     float64             `bson:"invested" json:"invested"`
	CurrentValue float64             `bson:"current_value" json:"current_value"`
	Returns      float64             `bson:"returns" json:"returns"` // Percentage
	Status       string              `bson:"status" json:"status"`   // Growing, Stable, Declining
	PurchaseDate time.Time           `bson:"purchase_date,omitempty" json:"purchase_date,omitempty"`
}

// CalculateReturns calculates the return percentage
func (i *Investment) CalculateReturns() {
	if i.Invested > 0 {
		i.Returns = ((i.CurrentValue - i.Invested) / i.Invested) * 100
	}
}

// UpdateStatus updates the status based on returns
func (i *Investment) UpdateStatus() {
	if i.Returns >= 10 {
		i.Status = "Growing"
	} else if i.Returns >= 0 {
		i.Status = "Stable"
	} else {
		i.Status = "Declining"
	}
}


