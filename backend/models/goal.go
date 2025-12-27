package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Goal struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	UserID        primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Name          string             `bson:"name" json:"name"`
	TargetAmount  float64            `bson:"target_amount" json:"target_amount"`
	CurrentAmount float64            `bson:"current_amount" json:"current_amount"`
	Deadline      time.Time          `bson:"deadline,omitempty" json:"deadline,omitempty"`
	Status        string             `bson:"status" json:"status"`     // Planned, In Progress, Completed
	Priority      string             `bson:"priority" json:"priority"` // High, Medium, Low
	Description   string             `bson:"description" json:"description"`
}

// CalculateProgress calculates the progress percentage
func (g *Goal) CalculateProgress() float64 {
	if g.TargetAmount > 0 {
		return (g.CurrentAmount / g.TargetAmount) * 100
	}
	return 0
}

// UpdateStatus updates the status based on progress
func (g *Goal) UpdateStatus() {
	progress := g.CalculateProgress()
	if progress >= 100 {
		g.Status = "Completed"
	} else if progress > 0 {
		g.Status = "In Progress"
	} else {
		g.Status = "Planned"
	}
}

// CalculateLinkedInvestmentsTotal calculates total current value of linked investments
func (g *Goal) CalculateLinkedInvestmentsTotal(investments []Investment) float64 {
	total := 0.0
	for _, inv := range investments {
		if inv.GoalID != nil && *inv.GoalID == g.ID {
			total += inv.CurrentValue
		}
	}
	return total
}
