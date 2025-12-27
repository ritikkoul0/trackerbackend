package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/middleware"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserFinancials retrieves user's financial data (income, expenses, savings)
func GetUserFinancials(c *gin.Context) {
	// Get email from context (set by auth middleware)
	email, exists := middleware.GetEmailFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserFinancials updates user's financial data
func UpdateUserFinancials(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
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

	// Find user
	var user models.User
	if err := config.DB.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided
	if updateData.MonthlyIncome != nil {
		user.MonthlyIncome = *updateData.MonthlyIncome
	}
	if updateData.MonthlyExpenses != nil {
		user.MonthlyExpenses = *updateData.MonthlyExpenses
	}
	if updateData.MonthlySavings != nil {
		user.MonthlySavings = *updateData.MonthlySavings
	}

	// Save updates
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

//
