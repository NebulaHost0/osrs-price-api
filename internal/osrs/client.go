package osrs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"osrs-price-api/internal/models"
)

const (
	// OSRS Wiki Real-time Prices API
	wikiAPIURL = "https://prices.runescape.wiki/api/v1/osrs/latest"
	// User-Agent is required by the OSRS Wiki API
	// Per https://oldschool.runescape.wiki/w/RuneScape:Real-time_Prices
	userAgent = "grandexchange.gg - OSRS GE Price Tracker"
)

// Client handles communication with OSRS data sources
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new OSRS client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetLatestPrices fetches the latest prices for all items
// Following OSRS Wiki API guidelines:
// - Uses bulk endpoint (not individual item requests)
// - Sets proper User-Agent
// - Respects rate limits (called max once per 5 minutes by worker)
func (c *Client) GetLatestPrices() (map[string]models.ItemPrice, error) {
	req, err := http.NewRequest("GET", wikiAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// OSRS Wiki API requires a User-Agent header
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var wikiResp models.OSRSWikiResponse
	if err := json.Unmarshal(body, &wikiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our internal format
	prices := make(map[string]models.ItemPrice)
	for idStr, data := range wikiResp.Data {
		prices[idStr] = models.ItemPrice{
			High:       data.High,
			HighTime:   data.HighTime,
			Low:        data.Low,
			LowTime:    data.LowTime,
			HighVolume: data.HighVolume,
			LowVolume:  data.LowVolume,
		}
	}

	return prices, nil
}

// GetItemPrice fetches the price for a specific item by ID
// Note: This correctly uses the bulk endpoint internally (GetLatestPrices)
// and extracts the single item, following OSRS Wiki best practices.
// We never make individual API calls per item to avoid hammering their servers.
func (c *Client) GetItemPrice(itemID string) (*models.ItemPrice, error) {
	prices, err := c.GetLatestPrices()
	if err != nil {
		return nil, err
	}

	price, exists := prices[itemID]
	if !exists {
		return nil, fmt.Errorf("item not found")
	}

	return &price, nil
}