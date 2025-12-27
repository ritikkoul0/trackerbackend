package controllers

import (
	"context"
	"investment-tracker-backend/config"
	"investment-tracker-backend/middleware"
	"investment-tracker-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserFinancials retrieves user's financial data (income, expenses, savings)
func GetUserFinancials(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get email from context (set by auth middleware)
	email, exists := middleware.GetEmailFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	collection := config.DB.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserFinancials updates user's financial data
func UpdateUserFinancials(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updateData struct {
		MonthlyIncome   *float64 `json:"monthly_income"`
		MonthlyExpenses *float64 `json:"monthly_expenses"`
		MonthlySavings  *float64 `json:"monthly_savings"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	setFields := update["$set"].(bson.M)

	if updateData.MonthlyIncome != nil {
		setFields["monthly_income"] = *updateData.MonthlyIncome
	}
	if updateData.MonthlyExpenses != nil {
		setFields["monthly_expenses"] = *updateData.MonthlyExpenses
	}
	if updateData.MonthlySavings != nil {
		setFields["monthly_savings"] = *updateData.MonthlySavings
	}

	collection := config.DB.Collection("users")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch and return updated user
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, user)
}


