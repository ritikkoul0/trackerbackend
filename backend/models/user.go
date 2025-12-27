package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	Email           string  `bson:"email" json:"email"`
	Name            string  `bson:"name,omitempty" json:"name,omitempty"`
	MonthlyIncome   float64 `bson:"monthly_income,omitempty" json:"monthly_income,omitempty"`
	MonthlyExpenses float64 `bson:"monthly_expenses,omitempty" json:"monthly_expenses,omitempty"`
	MonthlySavings  float64 `bson:"monthly_savings,omitempty" json:"monthly_savings,omitempty"`
}


