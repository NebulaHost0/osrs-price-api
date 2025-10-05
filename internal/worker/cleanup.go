package worker

import (
	"fmt"
	"log"
	"time"

	"osrs-price-api/internal/database"
)

// CleanupWorker handles database cleanup and aggregation
type CleanupWorker struct {
	repository *database.Repository
	interval   time.Duration
	stopChan   chan bool
}

// NewCleanupWorker creates a new cleanup worker
func NewCleanupWorker(repo *database.Repository, interval time.Duration) *CleanupWorker {
	return &CleanupWorker{
		repository: repo,
		interval:   interval,
		stopChan:   make(chan bool),
	}
}

// Start begins the periodic cleanup
func (cw *CleanupWorker) Start() {
	log.Printf("Starting cleanup worker (interval: %s)", cw.interval)

	// Run immediately on start
	cw.runCleanup()

	// Then continue on interval
	ticker := time.NewTicker(cw.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cw.runCleanup()
			case <-cw.stopChan:
				ticker.Stop()
				log.Println("Cleanup worker stopped")
				return
			}
		}
	}()
}

// Stop stops the cleanup worker
func (cw *CleanupWorker) Stop() {
	cw.stopChan <- true
}

func (cw *CleanupWorker) runCleanup() {
	log.Println("Running database cleanup and aggregation...")

	// Step 1: Aggregate data from 8-90 days ago into hourly buckets
	hourlyStart := time.Now().UTC().Add(-90 * 24 * time.Hour)
	hourlyEnd := time.Now().UTC().Add(-8 * 24 * time.Hour)
	
	hourlyAggregated, err := cw.repository.AggregateToHourly(hourlyStart, hourlyEnd)
	if err != nil {
		log.Printf("Error aggregating hourly data: %v", err)
	} else if hourlyAggregated > 0 {
		log.Printf("Aggregated %d records into hourly buckets", hourlyAggregated)
	}

	// Step 2: Aggregate hourly data older than 90 days into daily buckets
	dailyStart := time.Now().UTC().Add(-5 * 365 * 24 * time.Hour) // 5 years back
	dailyEnd := time.Now().UTC().Add(-90 * 24 * time.Hour)
	
	dailyAggregated, err := cw.repository.AggregateToDaily(dailyStart, dailyEnd)
	if err != nil {
		log.Printf("Error aggregating daily data: %v", err)
	} else if dailyAggregated > 0 {
		log.Printf("Aggregated %d records into daily buckets", dailyAggregated)
	}

	// Step 3: Delete raw 5-minute data older than 8 days (now that it's aggregated)
	rawDataCutoff := time.Now().UTC().Add(-8 * 24 * time.Hour)
	deleted, err := cw.repository.DeleteOldPriceHistory(rawDataCutoff)
	if err != nil {
		log.Printf("Error during cleanup: %v", err)
		return
	}

	if deleted > 0 {
		log.Printf("Deleted %d old raw records (older than 8 days)", deleted)
	}

	// Step 4: Delete hourly data older than 90 days (now aggregated to daily)
	hourlyDeleteCutoff := time.Now().UTC().Add(-90 * 24 * time.Hour)
	deletedHourly, err := cw.repository.DeleteOldHourlyData(hourlyDeleteCutoff)
	if err != nil {
		log.Printf("Error deleting old hourly data: %v", err)
	} else if deletedHourly > 0 {
		log.Printf("Deleted %d old hourly records (older than 90 days)", deletedHourly)
	}

	// Get database stats
	stats, err := cw.repository.GetDatabaseStats()
	if err != nil {
		log.Printf("Error getting database stats: %v", err)
	} else {
		log.Printf("Database stats - Raw: %d records, Size: %s, Oldest: %s",
			stats.TotalRecords,
			formatBytes(stats.EstimatedSize),
			stats.OldestRecord.Format("2006-01-02"))
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return "< 1 KB"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []rune{'K', 'M', 'G', 'T', 'P', 'E'}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), units[exp])
}