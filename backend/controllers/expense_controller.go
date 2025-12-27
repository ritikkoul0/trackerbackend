package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetExpenses retrieves all expenses
func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	if err := config.DB.Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetExpense retrieves a single expense by ID
func GetExpense(c *gin.Context) {
	id := c.Param("id")
	expenseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var expense models.Expense
	if err := config.DB.First(&expense, uint(expenseID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// CreateExpense creates a new expense
func CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	id := c.Param("id")
	expenseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get old expense to update budget
	var oldExpense models.Expense
	if err := config.DB.First(&oldExpense, uint(expenseID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ID = uint(expenseID)

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
	id := c.Param("id")
	expenseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var expense models.Expense
	if err := config.DB.First(&expense, uint(expenseID)).Error; err != nil {
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

// Made with Bob
