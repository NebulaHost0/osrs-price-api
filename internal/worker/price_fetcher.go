package worker

import (
	"log"
	"time"

	"osrs-price-api/internal/database"
	"osrs-price-api/internal/osrs"
)

// PriceFetcher periodically fetches and stores price data
type PriceFetcher struct {
	client     *osrs.Client
	repository *database.Repository
	interval   time.Duration
	stopChan   chan bool
}

// NewPriceFetcher creates a new price fetcher worker
func NewPriceFetcher(client *osrs.Client, repo *database.Repository, interval time.Duration) *PriceFetcher {
	return &PriceFetcher{
		client:     client,
		repository: repo,
		interval:   interval,
		stopChan:   make(chan bool),
	}
}

// Start begins the periodic price fetching
func (pf *PriceFetcher) Start() {
	log.Printf("Starting price fetcher worker (interval: %s)", pf.interval)
	
	// Fetch immediately on start
	pf.fetchAndStore()

	// Then continue on interval
	ticker := time.NewTicker(pf.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				pf.fetchAndStore()
			case <-pf.stopChan:
				ticker.Stop()
				log.Println("Price fetcher worker stopped")
				return
			}
		}
	}()
}

// Stop stops the price fetcher
func (pf *PriceFetcher) Stop() {
	pf.stopChan <- true
}

func (pf *PriceFetcher) fetchAndStore() {
	log.Println("Fetching latest prices from OSRS Wiki API...")
	
	prices, err := pf.client.GetLatestPrices()
	if err != nil {
		log.Printf("Error fetching prices: %v", err)
		return
	}

	log.Printf("Fetched %d item prices, saving to database...", len(prices))
	
	if err := pf.repository.SavePriceHistory(prices); err != nil {
		log.Printf("Error saving prices to database: %v", err)
		return
	}

	log.Println("Successfully saved prices to database")
}