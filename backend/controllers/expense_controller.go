package controllers

import (
	"context"
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetExpenses retrieves all expenses
func GetExpenses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.DB.Collection("expenses")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var expenses []models.Expense
	if err = cursor.All(ctx, &expenses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetExpense retrieves a single expense by ID
func GetExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.DB.Collection("expenses")
	var expense models.Expense

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&expense)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// CreateExpense creates a new expense
func CreateExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ID = primitive.NewObjectID()
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	collection := config.DB.Collection("expenses")
	_, err := collection.InsertOne(ctx, expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update budget total expenses if budget_id is provided
	if !expense.BudgetID.IsZero() {
		budgetCollection := config.DB.Collection("budgets")
		var budget models.Budget
		err := budgetCollection.FindOne(ctx, bson.M{"_id": expense.BudgetID}).Decode(&budget)
		if err == nil {
			budget.TotalExpenses += expense.Amount
			budget.CalculateSavings()
			budgetCollection.UpdateOne(ctx, bson.M{"_id": expense.BudgetID}, bson.M{"$set": budget})
		}
	}

	c.JSON(http.StatusCreated, expense)
}

// UpdateExpense updates an existing expense
func UpdateExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get old expense to update budget
	collection := config.DB.Collection("expenses")
	var oldExpense models.Expense
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&oldExpense)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ID = objectID
	expense.UpdatedAt = time.Now()

	update := bson.M{"$set": expense}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// Update budget total expenses
	if !oldExpense.BudgetID.IsZero() {
		budgetCollection := config.DB.Collection("budgets")
		var budget models.Budget
		err := budgetCollection.FindOne(ctx, bson.M{"_id": oldExpense.BudgetID}).Decode(&budget)
		if err == nil {
			budget.TotalExpenses -= oldExpense.Amount
			if expense.BudgetID == oldExpense.BudgetID {
				budget.TotalExpenses += expense.Amount
			}
			budget.CalculateSavings()
			budgetCollection.UpdateOne(ctx, bson.M{"_id": oldExpense.BudgetID}, bson.M{"$set": budget})
		}
	}

	c.JSON(http.StatusOK, expense)
}

// DeleteExpense deletes an expense
func DeleteExpense(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.DB.Collection("expenses")
	var expense models.Expense
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&expense)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// Update budget total expenses before deleting
	if !expense.BudgetID.IsZero() {
		budgetCollection := config.DB.Collection("budgets")
		var budget models.Budget
		err := budgetCollection.FindOne(ctx, bson.M{"_id": expense.BudgetID}).Decode(&budget)
		if err == nil {
			budget.TotalExpenses -= expense.Amount
			budget.CalculateSavings()
			budgetCollection.UpdateOne(ctx, bson.M{"_id": expense.BudgetID}, bson.M{"$set": budget})
		}
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}


