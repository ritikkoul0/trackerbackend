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

// GetInvestments retrieves all investments
func GetInvestments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.DB.Collection("investments")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var investments []models.Investment
	if err = cursor.All(ctx, &investments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// GetInvestment retrieves a single investment by ID
func GetInvestment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.DB.Collection("investments")
	var investment models.Investment

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&investment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// CreateInvestment creates a new investment
func CreateInvestment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var investment models.Investment
	if err := c.BindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
			"hint":  "Make sure all required fields are provided: name, type, invested, current_value",
		})
		return
	}

	// Manual validation
	if investment.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Investment name is required"})
		return
	}
	if investment.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Investment type is required"})
		return
	}
	if investment.Invested <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invested amount must be greater than 0"})
		return
	}

	investment.ID = primitive.NewObjectID()
	investment.CreatedAt = time.Now()
	investment.UpdatedAt = time.Now()

	// Set default user_id if not provided
	if investment.UserID.IsZero() {
		investment.UserID = primitive.NewObjectID()
	}

	// Set default purchase date if not provided
	if investment.PurchaseDate.IsZero() {
		investment.PurchaseDate = time.Now()
	}

	// Calculate returns and status
	investment.CalculateReturns()
	investment.UpdateStatus()

	collection := config.DB.Collection("investments")
	_, err := collection.InsertOne(ctx, investment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create investment: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, investment)
}

// UpdateInvestment updates an existing investment
func UpdateInvestment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the old investment to check if goal_id changed
	collection := config.DB.Collection("investments")
	var oldInvestment models.Investment
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&oldInvestment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	var investment models.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	investment.ID = objectID
	investment.UpdatedAt = time.Now()

	// Recalculate returns and status
	investment.CalculateReturns()
	investment.UpdateStatus()

	update := bson.M{"$set": investment}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Update goal current_amount if goal_id changed or current_value changed
	if oldInvestment.GoalID != nil {
		updateGoalCurrentAmount(ctx, *oldInvestment.GoalID)
	}
	if investment.GoalID != nil && (oldInvestment.GoalID == nil || *oldInvestment.GoalID != *investment.GoalID) {
		updateGoalCurrentAmount(ctx, *investment.GoalID)
	}

	c.JSON(http.StatusOK, investment)
}

// DeleteInvestment deletes an investment
func DeleteInvestment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := config.DB.Collection("investments")

	// Get the investment before deleting to check if it's linked to a goal
	var investment models.Investment
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&investment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Store the goal_id before deletion
	goalID := investment.GoalID

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Update the goal's current_amount if the investment was linked to a goal
	if goalID != nil {
		updateGoalCurrentAmount(ctx, *goalID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment deleted successfully"})
}

// LinkInvestmentToGoal links an investment to a goal
func LinkInvestmentToGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	investmentID := c.Param("id")
	investmentObjectID, err := primitive.ObjectIDFromHex(investmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid investment ID format"})
		return
	}

	var requestBody struct {
		GoalID string `json:"goal_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get the investment to check its previous goal_id
	collection := config.DB.Collection("investments")
	var investment models.Investment
	err = collection.FindOne(ctx, bson.M{"_id": investmentObjectID}).Decode(&investment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	var goalObjectID *primitive.ObjectID
	var newGoalID primitive.ObjectID
	if requestBody.GoalID != "" {
		gid, err := primitive.ObjectIDFromHex(requestBody.GoalID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID format"})
			return
		}
		goalObjectID = &gid
		newGoalID = gid

		// Verify goal exists
		goalCollection := config.DB.Collection("goals")
		var goal models.Goal
		err = goalCollection.FindOne(ctx, bson.M{"_id": gid}).Decode(&goal)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
			return
		}
	}

	// Update investment with goal_id
	update := bson.M{
		"$set": bson.M{
			"goal_id":    goalObjectID,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": investmentObjectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Update the old goal's current_amount (if there was one)
	if investment.GoalID != nil {
		updateGoalCurrentAmount(ctx, *investment.GoalID)
	}

	// Update the new goal's current_amount (if linking to a goal)
	if !newGoalID.IsZero() {
		updateGoalCurrentAmount(ctx, newGoalID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Investment linked to goal successfully",
		"goal_id": requestBody.GoalID,
	})
}

// updateGoalCurrentAmount recalculates and updates a goal's current_amount based on linked investments
func updateGoalCurrentAmount(ctx context.Context, goalID primitive.ObjectID) error {
	investmentCollection := config.DB.Collection("investments")
	goalCollection := config.DB.Collection("goals")

	// Find all investments linked to this goal
	cursor, err := investmentCollection.Find(ctx, bson.M{"goal_id": goalID})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var investments []models.Investment
	if err = cursor.All(ctx, &investments); err != nil {
		return err
	}

	// Calculate total current value
	totalCurrentValue := 0.0
	for _, inv := range investments {
		totalCurrentValue += inv.CurrentValue
	}

	// Update the goal's current_amount
	update := bson.M{
		"$set": bson.M{
			"current_amount": totalCurrentValue,
			"updated_at":     time.Now(),
		},
	}

	_, err = goalCollection.UpdateOne(ctx, bson.M{"_id": goalID}, update)
	return err
}

// UnlinkInvestmentFromGoal removes the goal link from an investment
func UnlinkInvestmentFromGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	investmentID := c.Param("id")
	investmentObjectID, err := primitive.ObjectIDFromHex(investmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid investment ID format"})
		return
	}

	// Get the investment to check its current goal_id
	collection := config.DB.Collection("investments")
	var investment models.Investment
	err = collection.FindOne(ctx, bson.M{"_id": investmentObjectID}).Decode(&investment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Store the old goal_id before unlinking
	oldGoalID := investment.GoalID

	// Update investment to remove goal_id
	update := bson.M{
		"$unset": bson.M{"goal_id": ""},
		"$set":   bson.M{"updated_at": time.Now()},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": investmentObjectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	// Update the old goal's current_amount
	if oldGoalID != nil {
		updateGoalCurrentAmount(ctx, *oldGoalID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment unlinked from goal successfully"})
}

// GetInvestmentsByGoal retrieves all investments linked to a specific goal
func GetInvestmentsByGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	goalID := c.Param("goal_id")
	goalObjectID, err := primitive.ObjectIDFromHex(goalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID format"})
		return
	}

	collection := config.DB.Collection("investments")
	cursor, err := collection.Find(ctx, bson.M{"goal_id": goalObjectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var investments []models.Investment
	if err = cursor.All(ctx, &investments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total
	total := 0.0
	for _, inv := range investments {
		total += inv.CurrentValue
	}

	c.JSON(http.StatusOK, gin.H{
		"investments": investments,
		"total":       total,
		"count":       len(investments),
	})
}
