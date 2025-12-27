package controllers

import (
	"context"
	"encoding/json"
	"investment-tracker-backend/config"
	"investment-tracker-backend/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  string
	jwtSecret         []byte
)

func InitOAuth() {
	// Load from environment variables
	oauthStateString = os.Getenv("OAUTH_STATE_STRING")
	if oauthStateString == "" {
		oauthStateString = "randomstate" // fallback
	}

	jwtSecretStr := os.Getenv("JWT_SECRET")
	if jwtSecretStr == "" {
		jwtSecretStr = "supersecretkey" // fallback
	}
	jwtSecret = []byte(jwtSecretStr)

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// HandleGoogleLogin initiates Google OAuth flow
func HandleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback handles the OAuth callback
func HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "code exchange failed"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	email := userInfo["email"].(string)
	name, _ := userInfo["name"].(string)

	// Save or update user in PostgreSQL
	user, err := saveOrUpdateUser(email, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	// Generate JWT token
	jwtToken, err := generateJWT(email, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Set cookie
	c.SetCookie("auth_token", jwtToken, 3600, "/", "localhost", false, true)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	c.Redirect(http.StatusFound, frontendURL+"/dashboard")
}

// saveOrUpdateUser saves or updates user in PostgreSQL
func saveOrUpdateUser(email, name string) (*models.User, error) {
	var user models.User

	// Try to find existing user
	result := config.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		// User doesn't exist, create new
		user = models.User{
			Email: email,
			Name:  name,
		}
		if err := config.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// User exists, update name
		user.Name = name
		if err := config.DB.Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// generateJWT creates a JWT token
func generateJWT(email, userID string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken verifies the JWT token
func VerifyToken(c *gin.Context) {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"email":   claims["email"],
			"user_id": claims["user_id"],
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
}

// GetUserInfo returns current user info
func GetUserInfo(c *gin.Context) {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"loggedIn": false})
		return
	}

	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"loggedIn": false})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.JSON(http.StatusOK, gin.H{
			"loggedIn": true,
			"email":    claims["email"],
			"user_id":  claims["user_id"],
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"loggedIn": false})
}

// Logout clears the auth cookie
func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Made with Bob
