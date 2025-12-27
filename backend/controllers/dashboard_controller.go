package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardResponse struct {
	TotalInvestments float64             `json:"total_investments"`
	TotalGains       float64             `json:"total_gains"`
	MonthlyIncome    float64             `json:"monthly_income"`
	MonthlyExpenses  float64             `json:"monthly_expenses"`
	MonthlySavings   float64             `json:"monthly_savings"`
	Investments      []models.Investment `json:"investments"`
	Goals            []models.Goal       `json:"goals"`
	RecentExpenses   []models.Expense    `json:"recent_expenses"`
}

// GetDashboard retrieves dashboard summary data for the authenticated user
func GetDashboard(c *gin.Context) {
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

	var response DashboardResponse

	// Get all investments for this user only
	var investments []models.Investment
	if err := config.DB.Where("user_id = ?", uint(userID)).Find(&investments).Error; err == nil {
		response.Investments = investments

		// Calculate total investments and gains
		var totalInvested float64
		var totalCurrent float64
		for _, inv := range investments {
			totalInvested += inv.Invested
			totalCurrent += inv.CurrentValue
		}
		response.TotalInvestments = totalCurrent
		response.TotalGains = totalCurrent - totalInvested
	}

	// Get current month budget (most recent) for this user only
	var budget models.Budget
	if err := config.DB.Where("user_id = ?", uint(userID)).Order("created_at DESC").First(&budget).Error; err == nil {
		response.MonthlyIncome = budget.Income
		response.MonthlyExpenses = budget.TotalExpenses
		response.MonthlySavings = budget.Savings
	}

	// Get active goals (not completed, limit 5) for this user only
	var goals []models.Goal
	if err := config.DB.Where("user_id = ? AND status != ?", uint(userID), "Completed").Limit(5).Find(&goals).Error; err == nil {
		response.Goals = goals
	}

	// Get recent expenses (limit 10, sorted by date descending) for this user only
	var expenses []models.Expense
	if err := config.DB.Where("user_id = ?", uint(userID)).Order("date DESC").Limit(10).Find(&expenses).Error; err == nil {
		response.RecentExpenses = expenses
	}

	c.JSON(http.StatusOK, response)
}

//
