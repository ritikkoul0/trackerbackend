package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetExpenses retrieves all expenses for the authenticated user
func GetExpenses(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var expenses []models.Expense
	if err := config.DB.Where("user_id = ?", uint(userID)).Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetExpense retrieves a single expense by ID for the authenticated user
func GetExpense(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	id := c.Param("id")
	expenseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var expense models.Expense
	if err := config.DB.Where("id = ? AND user_id = ?", uint(expenseID), uint(userID)).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// CreateExpense creates a new expense
func CreateExpense(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.UserID = uint(userID)

	if err := config.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update budget total expenses if budget_id is provided
	if expense.BudgetID != nil {
		var budget models.Budget
		if err := config.DB.First(&budget, *expense.BudgetID).Error; err == nil {
			budget.TotalExpenses += expense.Amount
			budget.CalculateSavings()
			config.DB.Save(&budget)
		}
	}

	c.JSON(http.StatusCreated, expense)
}

// UpdateExpense updates an existing expense
func UpdateExpense(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	id := c.Param("id")
	expenseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get old expense to update budget and verify ownership
	var oldExpense models.Expense
	if err := config.DB.Where("id = ? AND user_id = ?", uint(expenseID), uint(userID)).First(&oldExpense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ID = uint(expenseID)
	expense.UserID = uint(userID)

	if err := config.DB.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update budget total expenses
	if oldExpense.BudgetID != nil {
		var budget models.Budget
		if err := config.DB.First(&budget, *oldExpense.BudgetID).Error; err == nil {
			budget.TotalExpenses -= oldExpense.Amount
			if expense.BudgetID != nil && *expense.BudgetID == *oldExpense.BudgetID {
				budget.TotalExpenses += expense.Amount
			}
			budget.CalculateSavings()
			config.DB.Save(&budget)
		}
	}

	c.JSON(http.StatusOK, expense)
}

// DeleteExpense deletes an expense
func DeleteExpense(c *gin.Context) {
	// Get user_id from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	id := c.Param("id")
	expenseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Verify ownership before deleting
	var expense models.Expense
	if err := config.DB.Where("id = ? AND user_id = ?", uint(expenseID), uint(userID)).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// Update budget total expenses before deleting
	if expense.BudgetID != nil {
		var budget models.Budget
		if err := config.DB.First(&budget, *expense.BudgetID).Error; err == nil {
			budget.TotalExpenses -= expense.Amount
			budget.CalculateSavings()
			config.DB.Save(&budget)
		}
	}

	if err := config.DB.Delete(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

//
