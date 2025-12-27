package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"

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

// GetDashboard retrieves dashboard summary data
func GetDashboard(c *gin.Context) {
	var response DashboardResponse

	// Get all investments
	var investments []models.Investment
	if err := config.DB.Find(&investments).Error; err == nil {
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

	// Get current month budget (most recent)
	var budget models.Budget
	if err := config.DB.Order("created_at DESC").First(&budget).Error; err == nil {
		response.MonthlyIncome = budget.Income
		response.MonthlyExpenses = budget.TotalExpenses
		response.MonthlySavings = budget.Savings
	}

	// Get active goals (not completed, limit 5)
	var goals []models.Goal
	if err := config.DB.Where("status != ?", "Completed").Limit(5).Find(&goals).Error; err == nil {
		response.Goals = goals
	}

	// Get recent expenses (limit 10, sorted by date descending)
	var expenses []models.Expense
	if err := config.DB.Order("date DESC").Limit(10).Find(&expenses).Error; err == nil {
		response.RecentExpenses = expenses
	}

	c.JSON(http.StatusOK, response)
}

// Made with Bob
