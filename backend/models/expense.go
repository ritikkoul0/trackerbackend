package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	UserID      primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	BudgetID    primitive.ObjectID `bson:"budget_id,omitempty" json:"budget_id"`
	Category    string             `bson:"category" json:"category" binding:"required"` // Food, Transport, Entertainment, etc.
	Amount      float64            `bson:"amount" json:"amount" binding:"required"`
	Description string             `bson:"description" json:"description"`
	Date        time.Time          `bson:"date" json:"date"`
}


