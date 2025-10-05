package api

import (
	"net/http"
	"strconv"
	"time"

	"osrs-price-api/internal/cache"
	"osrs-price-api/internal/database"
	"osrs-price-api/internal/osrs"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	osrsClient *osrs.Client
	cache      *cache.PriceCache
	repository *database.Repository
}

// NewHandler creates a new API handler
func NewHandler(client *osrs.Client, cache *cache.PriceCache, repo *database.Repository) *Handler {
	return &Handler{
		osrsClient: client,
		cache:      cache,
		repository: repo,
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "osrs-price-api",
	})
}

// GetAllPrices returns all item prices
func (h *Handler) GetAllPrices(c *gin.Context) {
	// Check cache first
	if prices, found := h.cache.GetAll(); found {
		c.JSON(http.StatusOK, gin.H{
			"data":   prices,
			"cached": true,
		})
		return
	}

	// Fetch from API if not cached
	prices, err := h.osrsClient.GetLatestPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch prices",
			"message": err.Error(),
		})
		return
	}

	// Cache the results
	h.cache.SetAll(prices)

	c.JSON(http.StatusOK, gin.H{
		"data":   prices,
		"cached": false,
	})
}

// GetItemPrice returns the price for a specific item
func (h *Handler) GetItemPrice(c *gin.Context) {
	itemID := c.Param("id")

	// Validate item ID is numeric
	if _, err := strconv.Atoi(itemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID",
			"message": "Item ID must be a number",
		})
		return
	}

	// Check cache first
	if price, found := h.cache.Get(itemID); found {
		price.ID, _ = strconv.Atoi(itemID)
		c.JSON(http.StatusOK, gin.H{
			"data":   price,
			"cached": true,
		})
		return
	}

	// Fetch from API if not cached
	price, err := h.osrsClient.GetItemPrice(itemID)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found",
				"message": "No price data available for this item ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch price",
			"message": err.Error(),
		})
		return
	}

	// Cache the result
	h.cache.Set(itemID, *price)

	// Add ID to response
	price.ID, _ = strconv.Atoi(itemID)

	c.JSON(http.StatusOK, gin.H{
		"data":   price,
		"cached": false,
	})
}

// ClearCache clears all cached data
func (h *Handler) ClearCache(c *gin.Context) {
	h.cache.Clear()
	c.JSON(http.StatusOK, gin.H{
		"message": "Cache cleared successfully",
	})
}

// GetPriceHistory returns historical price data for an item
func (h *Handler) GetPriceHistory(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid item ID",
			"message": "Item ID must be a number",
		})
		return
	}

	// Parse time range parameters
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 24
	}

	endTime := time.Now().UTC()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	history, err := h.repository.GetPriceHistory(itemID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch price history",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item_id":    itemID,
		"start_time": startTime,
		"end_time":   endTime,
		"data":       history,
		"count":      len(history),
	})
}

// GetPriceChange returns price change data for an item
func (h *Handler) GetPriceChange(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid item ID",
			"message": "Item ID must be a number",
		})
		return
	}

	// Parse duration parameter (default 24 hours)
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 24
	}

	duration := time.Duration(hours) * time.Hour
	change, err := h.repository.GetPriceChange(itemID, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch price change",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": change,
	})
}

// GetPriceStats returns statistical data for an item
func (h *Handler) GetPriceStats(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid item ID",
			"message": "Item ID must be a number",
		})
		return
	}

	// Parse time range parameters
	hoursStr := c.DefaultQuery("hours", "168") // Default 7 days
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 168
	}

	endTime := time.Now().UTC()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	stats, err := h.repository.GetPriceStats(itemID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch price statistics",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// GetTopGainers returns items with the highest price increases
func (h *Handler) GetTopGainers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 24
	}

	duration := time.Duration(hours) * time.Hour
	gainers, err := h.repository.GetTopGainers(limit, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch top gainers",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       gainers,
		"count":      len(gainers),
		"time_range": duration.String(),
	})
}

// GetTopByVolume returns items with the highest trading volume
func (h *Handler) GetTopByVolume(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 24
	}

	duration := time.Duration(hours) * time.Hour
	items, err := h.repository.GetTopByVolume(limit, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch top items by volume",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       items,
		"count":      len(items),
		"time_range": duration.String(),
	})
}

// GetTopLosers returns items with the highest price decreases
func (h *Handler) GetTopLosers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 {
		hours = 24
	}

	duration := time.Duration(hours) * time.Hour
	losers, err := h.repository.GetTopLosers(limit, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch top losers",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       losers,
		"count":      len(losers),
		"time_range": duration.String(),
	})
}