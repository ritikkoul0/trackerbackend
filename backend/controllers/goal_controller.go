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

// GetGoals retrieves all goals
func GetGoals(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.DB.Collection("goals")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var goals []models.Goal
	if err = cursor.All(ctx, &goals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// GetGoal retrieves a single goal by ID
func GetGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.DB.Collection("goals")
	var goal models.Goal

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&goal)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// CreateGoal creates a new goal
func CreateGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	goal.ID = primitive.NewObjectID()
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	// Set default user_id if not provided
	if goal.UserID.IsZero() {
		goal.UserID = primitive.NewObjectID()
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

	collection := config.DB.Collection("goals")
	_, err := collection.InsertOne(ctx, goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// UpdateGoal updates an existing goal
func UpdateGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal.ID = objectID
	goal.UpdatedAt = time.Now()

	// Update status based on progress
	goal.UpdateStatus()

	collection := config.DB.Collection("goals")
	update := bson.M{"$set": goal}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoal deletes a goal
func DeleteGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.DB.Collection("goals")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}


