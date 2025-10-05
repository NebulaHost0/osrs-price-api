package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, handler *Handler) {
	// Health check
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Current prices
		v1.GET("/prices", handler.GetAllPrices)
		v1.GET("/prices/:id", handler.GetItemPrice)

		// Historical data
		v1.GET("/history/:id", handler.GetPriceHistory)
		v1.GET("/change/:id", handler.GetPriceChange)
		v1.GET("/stats/:id", handler.GetPriceStats)

		// Market analysis
		v1.GET("/gainers", handler.GetTopGainers)
		v1.GET("/losers", handler.GetTopLosers)
		v1.GET("/volume", handler.GetTopByVolume)

		// Cache management
		v1.POST("/cache/clear", handler.ClearCache)
	}
}
