package controllers

import (
	"context"
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var response DashboardResponse

	// Get all investments
	investmentCollection := config.DB.Collection("investments")
	cursor, err := investmentCollection.Find(ctx, bson.M{})
	if err == nil {
		var investments []models.Investment
		if err = cursor.All(ctx, &investments); err == nil {
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
		cursor.Close(ctx)
	}

	// Get current month budget (most recent)
	budgetCollection := config.DB.Collection("budgets")
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var budget models.Budget
	err = budgetCollection.FindOne(ctx, bson.M{}, opts).Decode(&budget)
	if err == nil {
		response.MonthlyIncome = budget.Income
		response.MonthlyExpenses = budget.TotalExpenses
		response.MonthlySavings = budget.Savings
	}

	// Get active goals (not completed, limit 5)
	goalCollection := config.DB.Collection("goals")
	goalOpts := options.Find().SetLimit(5)
	goalCursor, err := goalCollection.Find(ctx, bson.M{"status": bson.M{"$ne": "Completed"}}, goalOpts)
	if err == nil {
		var goals []models.Goal
		if err = goalCursor.All(ctx, &goals); err == nil {
			response.Goals = goals
		}
		goalCursor.Close(ctx)
	}

	// Get recent expenses (limit 10, sorted by date descending)
	expenseCollection := config.DB.Collection("expenses")
	expenseOpts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(10)
	expenseCursor, err := expenseCollection.Find(ctx, bson.M{}, expenseOpts)
	if err == nil {
		var expenses []models.Expense
		if err = expenseCursor.All(ctx, &expenses); err == nil {
			response.RecentExpenses = expenses
		}
		expenseCursor.Close(ctx)
	}

	c.JSON(http.StatusOK, response)
}


