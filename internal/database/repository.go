package database

import (
	"fmt"
	"time"

	"osrs-price-api/internal/models"

	"gorm.io/gorm"
)

// Repository handles database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// SavePriceHistory saves a batch of price history records
func (r *Repository) SavePriceHistory(prices map[string]models.ItemPrice) error {
	timestamp := time.Now().UTC()
	
	var records []models.PriceHistory
	for itemIDStr, price := range prices {
		var itemID int
		fmt.Sscanf(itemIDStr, "%d", &itemID)
		
		records = append(records, models.PriceHistory{
			ItemID:     itemID,
			High:       price.High,
			HighTime:   price.HighTime,
			Low:        price.Low,
			LowTime:    price.LowTime,
			HighVolume: price.HighVolume,
			LowVolume:  price.LowVolume,
			Timestamp:  timestamp,
		})
	}

	// Batch insert for better performance
	if len(records) > 0 {
		return r.db.CreateInBatches(records, 100).Error
	}

	return nil
}

// GetLatestPrice retrieves the most recent price for an item
func (r *Repository) GetLatestPrice(itemID int) (*models.PriceHistory, error) {
	var price models.PriceHistory
	err := r.db.Where("item_id = ?", itemID).
		Order("timestamp DESC").
		First(&price).Error
	
	if err != nil {
		return nil, err
	}
	return &price, nil
}

// GetPriceHistory retrieves price history for an item within a time range
// Automatically uses the appropriate table based on time range:
// - Raw data (5-min) for last 7 days
// - Hourly aggregates for 8-90 days
// - Daily aggregates for 90+ days
func (r *Repository) GetPriceHistory(itemID int, startTime, endTime time.Time) ([]models.PriceHistory, error) {
	now := time.Now().UTC()
	daysAgo := int(now.Sub(startTime).Hours() / 24)
	
	// If requesting recent data (within 7 days), use raw 5-minute data
	if daysAgo <= 7 {
		var history []models.PriceHistory
		err := r.db.Where("item_id = ? AND timestamp BETWEEN ? AND ?", itemID, startTime, endTime).
			Order("timestamp ASC").
			Find(&history).Error
		return history, err
	}
	
	// If requesting 8-90 days, use hourly aggregates
	if daysAgo <= 90 {
		var hourlyData []models.PriceHistoryHourly
		err := r.db.Where("item_id = ? AND hour_timestamp BETWEEN ? AND ?", itemID, startTime, endTime).
			Order("hour_timestamp ASC").
			Find(&hourlyData).Error
		
		// Convert to PriceHistory format
		history := make([]models.PriceHistory, len(hourlyData))
		for i, h := range hourlyData {
			history[i] = models.PriceHistory{
				ItemID:     h.ItemID,
				High:       h.AvgHigh,
				Low:        h.AvgLow,
				HighVolume: h.TotalHighVolume,
				LowVolume:  h.TotalLowVolume,
				Timestamp:  h.HourTimestamp,
			}
		}
		return history, err
	}
	
	// For 90+ days, use daily aggregates
	var dailyData []models.PriceHistoryDaily
	err := r.db.Where("item_id = ? AND day_date BETWEEN ? AND ?", itemID, startTime.Truncate(24*time.Hour), endTime.Truncate(24*time.Hour)).
		Order("day_date ASC").
		Find(&dailyData).Error
	
	// Convert to PriceHistory format
	history := make([]models.PriceHistory, len(dailyData))
	for i, d := range dailyData {
		history[i] = models.PriceHistory{
			ItemID:     d.ItemID,
			High:       d.AvgHigh,
			Low:        d.AvgLow,
			HighVolume: d.TotalHighVolume,
			LowVolume:  d.TotalLowVolume,
			Timestamp:  d.DayDate,
		}
	}
	return history, err
}

// GetPriceChange calculates price change for an item over a time period
func (r *Repository) GetPriceChange(itemID int, duration time.Duration) (*models.PriceChangeResponse, error) {
	now := time.Now().UTC()
	startTime := now.Add(-duration)

	var current, previous models.PriceHistory

	// Get current (most recent) price
	if err := r.db.Where("item_id = ?", itemID).
		Order("timestamp DESC").
		First(&current).Error; err != nil {
		return nil, fmt.Errorf("no current price data: %w", err)
	}

	// Get previous price (closest to start time)
	if err := r.db.Where("item_id = ? AND timestamp >= ?", itemID, startTime).
		Order("timestamp ASC").
		First(&previous).Error; err != nil {
		return nil, fmt.Errorf("no historical price data: %w", err)
	}

	highChange := current.High - previous.High
	lowChange := current.Low - previous.Low

	var highChangePerc, lowChangePerc float64
	if previous.High > 0 {
		highChangePerc = (float64(highChange) / float64(previous.High)) * 100
	}
	if previous.Low > 0 {
		lowChangePerc = (float64(lowChange) / float64(previous.Low)) * 100
	}

	return &models.PriceChangeResponse{
		ItemID:         itemID,
		CurrentHigh:    current.High,
		CurrentLow:     current.Low,
		PreviousHigh:   previous.High,
		PreviousLow:    previous.Low,
		HighChange:     highChange,
		LowChange:      lowChange,
		HighChangePerc: highChangePerc,
		LowChangePerc:  lowChangePerc,
		TimeRange:      duration.String(),
		Timestamp:      current.Timestamp,
	}, nil
}

// GetPriceStats calculates statistical data for an item
func (r *Repository) GetPriceStats(itemID int, startTime, endTime time.Time) (*models.PriceStats, error) {
	var stats struct {
		AvgHigh    float64
		AvgLow     float64
		MaxHigh    int64
		MaxLow     int64
		MinHigh    int64
		MinLow     int64
		DataPoints int64
	}

	err := r.db.Model(&models.PriceHistory{}).
		Select(`
			AVG(high) as avg_high,
			AVG(low) as avg_low,
			MAX(high) as max_high,
			MAX(low) as max_low,
			MIN(high) as min_high,
			MIN(low) as min_low,
			COUNT(*) as data_points
		`).
		Where("item_id = ? AND timestamp BETWEEN ? AND ?", itemID, startTime, endTime).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// Calculate volatility (standard deviation of high prices)
	var volatility float64
	r.db.Model(&models.PriceHistory{}).
		Select("STDDEV(high)").
		Where("item_id = ? AND timestamp BETWEEN ? AND ?", itemID, startTime, endTime).
		Scan(&volatility)

	return &models.PriceStats{
		ItemID:     itemID,
		AvgHigh:    stats.AvgHigh,
		AvgLow:     stats.AvgLow,
		MaxHigh:    stats.MaxHigh,
		MaxLow:     stats.MaxLow,
		MinHigh:    stats.MinHigh,
		MinLow:     stats.MinLow,
		Volatility: volatility,
		DataPoints: stats.DataPoints,
		StartTime:  startTime,
		EndTime:    endTime,
	}, nil
}

// GetTopGainers returns items with the highest price increases
func (r *Repository) GetTopGainers(limit int, duration time.Duration) ([]models.PriceChangeResponse, error) {
	now := time.Now().UTC()
	startTime := now.Add(-duration)

	query := `
		WITH current_prices AS (
			SELECT DISTINCT ON (item_id)
				item_id, high, low, timestamp
			FROM price_history
			WHERE timestamp >= $1
			ORDER BY item_id, timestamp DESC
		),
		previous_prices AS (
			SELECT DISTINCT ON (item_id)
				item_id, high as prev_high, low as prev_low
			FROM price_history
			WHERE timestamp >= $2
			ORDER BY item_id, timestamp ASC
		)
		SELECT 
			c.item_id,
			c.high as current_high,
			c.low as current_low,
			p.prev_high as previous_high,
			p.prev_low as previous_low,
			(c.high - p.prev_high) as high_change,
			(c.low - p.prev_low) as low_change,
			CASE WHEN p.prev_high > 0 THEN ((c.high - p.prev_high)::float / p.prev_high) * 100 ELSE 0 END as high_change_perc,
			CASE WHEN p.prev_low > 0 THEN ((c.low - p.prev_low)::float / p.prev_low) * 100 ELSE 0 END as low_change_perc,
			c.timestamp
		FROM current_prices c
		JOIN previous_prices p ON c.item_id = p.item_id
		WHERE p.prev_high > 0 AND c.high > p.prev_high
		ORDER BY high_change_perc DESC
		LIMIT $3
	`

	var results []struct {
		ItemID         int
		CurrentHigh    int64
		CurrentLow     int64
		PreviousHigh   int64
		PreviousLow    int64
		HighChange     int64
		LowChange      int64
		HighChangePerc float64
		LowChangePerc  float64
		Timestamp      time.Time
	}

	err := r.db.Raw(query, now, startTime, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var changes []models.PriceChangeResponse
	for _, r := range results {
		changes = append(changes, models.PriceChangeResponse{
			ItemID:         r.ItemID,
			CurrentHigh:    r.CurrentHigh,
			CurrentLow:     r.CurrentLow,
			PreviousHigh:   r.PreviousHigh,
			PreviousLow:    r.PreviousLow,
			HighChange:     r.HighChange,
			LowChange:      r.LowChange,
			HighChangePerc: r.HighChangePerc,
			LowChangePerc:  r.LowChangePerc,
			TimeRange:      duration.String(),
			Timestamp:      r.Timestamp,
		})
	}

	return changes, nil
}

// GetTopByVolume returns items with the highest trading volume
func (r *Repository) GetTopByVolume(limit int, duration time.Duration) ([]struct {
	ItemID      int   `json:"item_id"`
	TotalVolume int64 `json:"total_volume"`
	AvgHigh     int64 `json:"avg_high"`
	AvgLow      int64 `json:"avg_low"`
}, error) {
	cutoff := time.Now().UTC().Add(-duration)
	
	var results []struct {
		ItemID      int   `json:"item_id"`
		TotalVolume int64 `json:"total_volume"`
		AvgHigh     int64 `json:"avg_high"`
		AvgLow      int64 `json:"avg_low"`
	}
	
	err := r.db.Model(&models.PriceHistory{}).
		Select("item_id, SUM(high_volume + low_volume) as total_volume, AVG(high) as avg_high, AVG(low) as avg_low").
		Where("timestamp > ?", cutoff).
		Group("item_id").
		Having("SUM(high_volume + low_volume) > 0").
		Order("total_volume DESC").
		Limit(limit).
		Scan(&results).Error
	
	return results, err
}

// AggregateToHourly aggregates 5-minute data into hourly buckets
func (r *Repository) AggregateToHourly(startTime, endTime time.Time) (int64, error) {
	// Aggregate data for each hour in the range
	query := `
		INSERT INTO price_history_hourly (
			item_id, avg_high, avg_low, max_high, min_low,
			opening_high, opening_low, closing_high, closing_low,
			total_high_volume, total_low_volume, data_points, hour_timestamp
		)
		SELECT 
			item_id,
			AVG(high) as avg_high,
			AVG(low) as avg_low,
			MAX(high) as max_high,
			MIN(low) as min_low,
			FIRST_VALUE(high) OVER (PARTITION BY item_id, date_trunc('hour', timestamp) ORDER BY timestamp) as opening_high,
			FIRST_VALUE(low) OVER (PARTITION BY item_id, date_trunc('hour', timestamp) ORDER BY timestamp) as opening_low,
			LAST_VALUE(high) OVER (PARTITION BY item_id, date_trunc('hour', timestamp) ORDER BY timestamp ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) as closing_high,
			LAST_VALUE(low) OVER (PARTITION BY item_id, date_trunc('hour', timestamp) ORDER BY timestamp ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) as closing_low,
			SUM(high_volume) as total_high_volume,
			SUM(low_volume) as total_low_volume,
			COUNT(*) as data_points,
			date_trunc('hour', timestamp) as hour_timestamp
		FROM price_history
		WHERE timestamp >= ? AND timestamp < ?
		GROUP BY item_id, date_trunc('hour', timestamp)
		ON CONFLICT DO NOTHING
	`
	
	result := r.db.Exec(query, startTime, endTime)
	return result.RowsAffected, result.Error
}

// AggregateToDaily aggregates hourly data into daily buckets
func (r *Repository) AggregateToDaily(startTime, endTime time.Time) (int64, error) {
	query := `
		INSERT INTO price_history_daily (
			item_id, avg_high, avg_low, max_high, min_low,
			opening_high, opening_low, closing_high, closing_low,
			total_high_volume, total_low_volume, volatility, data_points, day_date
		)
		SELECT 
			item_id,
			AVG(avg_high) as avg_high,
			AVG(avg_low) as avg_low,
			MAX(max_high) as max_high,
			MIN(min_low) as min_low,
			FIRST_VALUE(opening_high) OVER (PARTITION BY item_id, DATE(hour_timestamp) ORDER BY hour_timestamp) as opening_high,
			FIRST_VALUE(opening_low) OVER (PARTITION BY item_id, DATE(hour_timestamp) ORDER BY hour_timestamp) as opening_low,
			LAST_VALUE(closing_high) OVER (PARTITION BY item_id, DATE(hour_timestamp) ORDER BY hour_timestamp ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) as closing_high,
			LAST_VALUE(closing_low) OVER (PARTITION BY item_id, DATE(hour_timestamp) ORDER BY hour_timestamp ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) as closing_low,
			SUM(total_high_volume) as total_high_volume,
			SUM(total_low_volume) as total_low_volume,
			CASE WHEN AVG(avg_high) > 0 THEN (MAX(max_high) - MIN(min_low))::float / AVG(avg_high) ELSE 0 END as volatility,
			SUM(data_points) as data_points,
			DATE(hour_timestamp) as day_date
		FROM price_history_hourly
		WHERE hour_timestamp >= ? AND hour_timestamp < ?
		GROUP BY item_id, DATE(hour_timestamp)
		ON CONFLICT (item_id, day_date) DO NOTHING
	`
	
	result := r.db.Exec(query, startTime, endTime)
	return result.RowsAffected, result.Error
}

// DeleteOldPriceHistory deletes price history older than the given date
func (r *Repository) DeleteOldPriceHistory(cutoffDate time.Time) (int64, error) {
	result := r.db.Where("timestamp < ?", cutoffDate).Delete(&models.PriceHistory{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// DeleteOldHourlyData deletes hourly aggregates older than the given date
func (r *Repository) DeleteOldHourlyData(cutoffDate time.Time) (int64, error) {
	result := r.db.Exec("DELETE FROM price_history_hourly WHERE hour_timestamp < ?", cutoffDate)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// DatabaseStats represents database statistics
type DatabaseStats struct {
	TotalRecords  int64
	OldestRecord  time.Time
	NewestRecord  time.Time
	EstimatedSize int64
}

// GetDatabaseStats returns statistics about the price_history table
func (r *Repository) GetDatabaseStats() (*DatabaseStats, error) {
	var stats DatabaseStats
	
	// Get total record count
	if err := r.db.Model(&models.PriceHistory{}).Count(&stats.TotalRecords).Error; err != nil {
		return nil, err
	}
	
	// Get oldest record
	var oldest models.PriceHistory
	if err := r.db.Order("timestamp ASC").First(&oldest).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	stats.OldestRecord = oldest.Timestamp
	
	// Get newest record
	var newest models.PriceHistory
	if err := r.db.Order("timestamp DESC").First(&newest).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	stats.NewestRecord = newest.Timestamp
	
	// Estimate size (96 bytes per record is approximate)
	stats.EstimatedSize = stats.TotalRecords * 96
	
	return &stats, nil
}

// GetTopLosers returns items with the highest price decreases
func (r *Repository) GetTopLosers(limit int, duration time.Duration) ([]models.PriceChangeResponse, error) {
	now := time.Now().UTC()
	startTime := now.Add(-duration)

	query := `
		WITH current_prices AS (
			SELECT DISTINCT ON (item_id)
				item_id, high, low, timestamp
			FROM price_history
			WHERE timestamp >= $1
			ORDER BY item_id, timestamp DESC
		),
		previous_prices AS (
			SELECT DISTINCT ON (item_id)
				item_id, high as prev_high, low as prev_low
			FROM price_history
			WHERE timestamp >= $2
			ORDER BY item_id, timestamp ASC
		)
		SELECT 
			c.item_id,
			c.high as current_high,
			c.low as current_low,
			p.prev_high as previous_high,
			p.prev_low as previous_low,
			(c.high - p.prev_high) as high_change,
			(c.low - p.prev_low) as low_change,
			CASE WHEN p.prev_high > 0 THEN ((c.high - p.prev_high)::float / p.prev_high) * 100 ELSE 0 END as high_change_perc,
			CASE WHEN p.prev_low > 0 THEN ((c.low - p.prev_low)::float / p.prev_low) * 100 ELSE 0 END as low_change_perc,
			c.timestamp
		FROM current_prices c
		JOIN previous_prices p ON c.item_id = p.item_id
		WHERE p.prev_high > 0 AND c.high < p.prev_high
		ORDER BY high_change_perc ASC
		LIMIT $3
	`

	var results []struct {
		ItemID         int
		CurrentHigh    int64
		CurrentLow     int64
		PreviousHigh   int64
		PreviousLow    int64
		HighChange     int64
		LowChange      int64
		HighChangePerc float64
		LowChangePerc  float64
		Timestamp      time.Time
	}

	err := r.db.Raw(query, now, startTime, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var changes []models.PriceChangeResponse
	for _, r := range results {
		changes = append(changes, models.PriceChangeResponse{
			ItemID:         r.ItemID,
			CurrentHigh:    r.CurrentHigh,
			CurrentLow:     r.CurrentLow,
			PreviousHigh:   r.PreviousHigh,
			PreviousLow:    r.PreviousLow,
			HighChange:     r.HighChange,
			LowChange:      r.LowChange,
			HighChangePerc: r.HighChangePerc,
			LowChangePerc:  r.LowChangePerc,
			TimeRange:      duration.String(),
			Timestamp:      r.Timestamp,
		})
	}

	return changes, nil
}