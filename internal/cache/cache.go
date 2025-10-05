package cache

import (
	"time"

	"osrs-price-api/internal/models"

	gocache "github.com/patrickmn/go-cache"
)

const (
	// Cache prices for 5 minutes (OSRS Wiki updates every 5 minutes)
	defaultExpiration = 5 * time.Minute
	// Cleanup expired items every 10 minutes
	cleanupInterval = 10 * time.Minute
)

// PriceCache handles caching of price data
type PriceCache struct {
	cache *gocache.Cache
}

// NewPriceCache creates a new price cache
func NewPriceCache() *PriceCache {
	return &PriceCache{
		cache: gocache.New(defaultExpiration, cleanupInterval),
	}
}

// Get retrieves a cached item price
func (pc *PriceCache) Get(itemID string) (*models.ItemPrice, bool) {
	if val, found := pc.cache.Get(itemID); found {
		if price, ok := val.(models.ItemPrice); ok {
			return &price, true
		}
	}
	return nil, false
}

// Set stores an item price in the cache
func (pc *PriceCache) Set(itemID string, price models.ItemPrice) {
	pc.cache.Set(itemID, price, defaultExpiration)
}

// GetAll retrieves all cached prices
func (pc *PriceCache) GetAll() (map[string]models.ItemPrice, bool) {
	if val, found := pc.cache.Get("all_prices"); found {
		if prices, ok := val.(map[string]models.ItemPrice); ok {
			return prices, true
		}
	}
	return nil, false
}

// SetAll stores all prices in the cache
func (pc *PriceCache) SetAll(prices map[string]models.ItemPrice) {
	pc.cache.Set("all_prices", prices, defaultExpiration)
}

// Clear removes all items from the cache
func (pc *PriceCache) Clear() {
	pc.cache.Flush()
}