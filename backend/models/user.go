package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email           string  `gorm:"uniqueIndex;not null" json:"email"`
	Name            string  `gorm:"type:varchar(255)" json:"name,omitempty"`
	MonthlyIncome   float64 `gorm:"type:decimal(15,2);default:0" json:"monthly_income,omitempty"`
	MonthlyExpenses float64 `gorm:"type:decimal(15,2);default:0" json:"monthly_expenses,omitempty"`
	MonthlySavings  float64 `gorm:"type:decimal(15,2);default:0" json:"monthly_savings,omitempty"`
}

//
