package models

import (
	"time"
)

// PriceHistory stores historical price data for items
type PriceHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ItemID         int       `gorm:"index:idx_item_timestamp;not null" json:"item_id"`
	High           int64     `json:"high"`
	HighTime       int64     `json:"high_time"`
	Low            int64     `json:"low"`
	LowTime        int64     `json:"low_time"`
	HighVolume     int64     `json:"high_volume"`     // Volume of high price trades
	LowVolume      int64     `json:"low_volume"`      // Volume of low price trades
	Timestamp      time.Time `gorm:"index:idx_item_timestamp;not null" json:"timestamp"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (PriceHistory) TableName() string {
	return "price_history"
}

// PriceChangeResponse represents price change data over time
type PriceChangeResponse struct {
	ItemID         int       `json:"item_id"`
	CurrentHigh    int64     `json:"current_high"`
	CurrentLow     int64     `json:"current_low"`
	PreviousHigh   int64     `json:"previous_high"`
	PreviousLow    int64     `json:"previous_low"`
	HighChange     int64     `json:"high_change"`
	LowChange      int64     `json:"low_change"`
	HighChangePerc float64   `json:"high_change_percent"`
	LowChangePerc  float64   `json:"low_change_percent"`
	TimeRange      string    `json:"time_range"`
	Timestamp      time.Time `json:"timestamp"`
}

// PriceStats represents statistical data for an item
type PriceStats struct {
	ItemID      int     `json:"item_id"`
	AvgHigh     float64 `json:"avg_high"`
	AvgLow      float64 `json:"avg_low"`
	MaxHigh     int64   `json:"max_high"`
	MaxLow      int64   `json:"max_low"`
	MinHigh     int64   `json:"min_high"`
	MinLow      int64   `json:"min_low"`
	Volatility  float64 `json:"volatility"`
	DataPoints  int64   `json:"data_points"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}