package router

import (
	"fmt"
	"github.com/khedhrije/podcaster-indexer-api/internal/configuration"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	validationEndpoint = "/private/token/validate"
	refreshEndpoint    = "/private/token/refresh"
)

// TokenValidatorMiddleware creates a Gin middleware that validates a token by calling an external endpoint.
func TokenValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Split the header to get the token part
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		url := fmt.Sprintf("%s%s", configuration.Config.AccountApi.BaseURL, validationEndpoint)
		// Create a new HTTP request
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Read the response body
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			c.Abort()
			return
		}

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			TokenRefresherMiddleware()
			return
		}

		// Token is valid, proceed to the next handler
		c.Next()
	}
}

// TokenRefresherMiddleware creates a Gin middleware that refreshes a token by calling an external endpoint.
func TokenRefresherMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		refreshHeader := c.GetHeader("refresh_token")
		if refreshHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh header is required"})
			c.Abort()
			return
		}

		// Split the header to get the token part
		tokenParts := strings.Split(refreshHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Create a new HTTP request
		url := fmt.Sprintf("%s%s", configuration.Config.AccountApi.BaseURL, refreshEndpoint)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("refresh_token", "Bearer "+token)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			c.Abort()
			return
		}

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token", "details": string(body)})
			c.Abort()
			return
		}

		// Parse the new token from the response
		newToken := string(body) // Assuming the body contains the new token as plain text, adjust as needed

		// Set the new token in the refresh_token header for the response
		c.Header("refresh_token", "Bearer "+newToken)

		// Proceed to the next handler
		c.Next()
	}
}
