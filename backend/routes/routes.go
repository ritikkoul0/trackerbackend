package routes

import (
	"investment-tracker-backend/controllers"
	"investment-tracker-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Auth routes (public - no middleware)
	router.GET("/sso", controllers.HandleGoogleLogin)
	router.GET("/auth/google/callback", controllers.HandleGoogleCallback)
	router.GET("/verify", controllers.VerifyToken)
	router.POST("/logout", controllers.Logout)
	router.GET("/me", controllers.GetUserInfo)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Health check (public)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
		})

		// Protected routes - require authentication
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Investment routes
			investments := protected.Group("/investments")
			{
				investments.GET("", controllers.GetInvestments)
				investments.GET("/:id", controllers.GetInvestment)
				investments.POST("", controllers.CreateInvestment)
				investments.PUT("/:id", controllers.UpdateInvestment)
				investments.DELETE("/:id", controllers.DeleteInvestment)
				investments.POST("/:id/link-goal", controllers.LinkInvestmentToGoal)
				investments.POST("/:id/unlink-goal", controllers.UnlinkInvestmentFromGoal)
				investments.GET("/by-goal/:goal_id", controllers.GetInvestmentsByGoal)
			}

			// Goal routes
			goals := protected.Group("/goals")
			{
				goals.GET("", controllers.GetGoals)
				goals.GET("/:id", controllers.GetGoal)
				goals.POST("", controllers.CreateGoal)
				goals.PUT("/:id", controllers.UpdateGoal)
				goals.DELETE("/:id", controllers.DeleteGoal)
			}

			// Budget routes
			budgets := protected.Group("/budgets")
			{
				budgets.GET("", controllers.GetBudgets)
				budgets.GET("/:id", controllers.GetBudget)
				budgets.POST("", controllers.CreateBudget)
				budgets.PUT("/:id", controllers.UpdateBudget)
				budgets.DELETE("/:id", controllers.DeleteBudget)
			}

			// Expense routes
			expenses := protected.Group("/expenses")
			{
				expenses.GET("", controllers.GetExpenses)
				expenses.GET("/:id", controllers.GetExpense)
				expenses.POST("", controllers.CreateExpense)
				expenses.PUT("/:id", controllers.UpdateExpense)
				expenses.DELETE("/:id", controllers.DeleteExpense)
			}

			// User routes
			users := protected.Group("/users")
			{
				users.GET("/financials", controllers.GetUserFinancials)
				users.PUT("/:id/financials", controllers.UpdateUserFinancials)
			}

			// Dashboard route
			protected.GET("/dashboard", controllers.GetDashboard)
		}
	}
}
