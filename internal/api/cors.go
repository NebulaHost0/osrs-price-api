package api

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Default allowed origins for development and production
		allowedOrigins := []string{
			// Local development
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
			"https://localhost:3000",
			// Production domains
			"https://grandexchange.gg",
			"https://www.grandexchange.gg",
			"http://grandexchange.gg",
			"http://www.grandexchange.gg",
		}

		// Allow additional origins from environment variable
		extraOrigins := os.Getenv("ALLOWED_ORIGINS")
		if extraOrigins != "" {
			origins := strings.Split(extraOrigins, ",")
			for _, o := range origins {
				allowedOrigins = append(allowedOrigins, strings.TrimSpace(o))
			}
		}

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		// Set CORS headers
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin == "" {
			// Allow requests without Origin header (curl, server-to-server, etc.)
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Log rejected origins for debugging
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Standard CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}