package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	BudgetID    *uint     `gorm:"index" json:"budget_id,omitempty"`
	Budget      *Budget   `gorm:"foreignKey:BudgetID;constraint:OnDelete:SET NULL" json:"-"`
	Category    string    `gorm:"type:varchar(100);not null" json:"category" binding:"required"` // Food, Transport, Entertainment, etc.
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
}

//
