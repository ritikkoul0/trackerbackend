package models

import (
	"time"

	"gorm.io/gorm"
)

type Investment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID       uint      `gorm:"not null;index" json:"user_id,omitempty"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	GoalID       *uint     `gorm:"index" json:"goal_id,omitempty"` // Link to Goal
	Goal         *Goal     `gorm:"foreignKey:GoalID;constraint:OnDelete:SET NULL" json:"-"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Type         string    `gorm:"type:varchar(100);not null" json:"type"` // Stocks, Mutual Fund, ETF, FD, PPF, etc.
	Invested     float64   `gorm:"type:decimal(15,2);not null" json:"invested"`
	CurrentValue float64   `gorm:"type:decimal(15,2);default:0" json:"current_value"`
	Returns      float64   `gorm:"type:decimal(10,2);default:0" json:"returns"`     // Percentage
	Status       string    `gorm:"type:varchar(50);default:'Stable'" json:"status"` // Growing, Stable, Declining
	PurchaseDate time.Time `json:"purchase_date,omitempty"`
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

//
