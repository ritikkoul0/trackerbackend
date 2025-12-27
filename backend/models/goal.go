package models

import (
	"time"

	"gorm.io/gorm"
)

type Goal struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID        uint      `gorm:"not null;index" json:"user_id,omitempty"`
	User          User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	TargetAmount  float64   `gorm:"type:decimal(15,2);not null" json:"target_amount"`
	CurrentAmount float64   `gorm:"type:decimal(15,2);default:0" json:"current_amount"`
	Deadline      time.Time `json:"deadline,omitempty"`
	Status        string    `gorm:"type:varchar(50);default:'Planned'" json:"status"`  // Planned, In Progress, Completed
	Priority      string    `gorm:"type:varchar(50);default:'Medium'" json:"priority"` // High, Medium, Low
	Description   string    `gorm:"type:text" json:"description"`
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

// Made with Bob
