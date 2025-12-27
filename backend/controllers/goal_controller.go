package controllers

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetGoals retrieves all goals
func GetGoals(c *gin.Context) {
	var goals []models.Goal
	if err := config.DB.Find(&goals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetGoal retrieves a single goal by ID
func GetGoal(c *gin.Context) {
	id := c.Param("id")
	goalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var goal models.Goal
	if err := config.DB.First(&goal, uint(goalID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// CreateGoal creates a new goal
func CreateGoal(c *gin.Context) {
	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Manual validation
	if goal.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Goal name is required"})
		return
	}
	if goal.TargetAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target amount must be greater than 0"})
		return
	}

	// Set default deadline if not provided
	if goal.Deadline.IsZero() {
		goal.Deadline = time.Now().AddDate(1, 0, 0) // 1 year from now
	}

	// Set default status if not provided
	if goal.Status == "" {
		goal.Status = "Planned"
	}

	// Update status based on progress
	goal.UpdateStatus()

	if err := config.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// UpdateGoal updates an existing goal
func UpdateGoal(c *gin.Context) {
	id := c.Param("id")
	goalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal.ID = uint(goalID)

	// Update status based on progress
	goal.UpdateStatus()

	if err := config.DB.Save(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoal deletes a goal
func DeleteGoal(c *gin.Context) {
	id := c.Param("id")
	goalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := config.DB.Delete(&models.Goal{}, uint(goalID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}

// Made with Bob
