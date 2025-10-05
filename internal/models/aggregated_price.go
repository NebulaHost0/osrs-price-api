package models

import "time"

// PriceHistoryHourly stores hourly aggregated price data
type PriceHistoryHourly struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ItemID          int       `gorm:"index:idx_hourly_item_time;not null" json:"item_id"`
	AvgHigh         int64     `json:"avg_high"`
	AvgLow          int64     `json:"avg_low"`
	MaxHigh         int64     `json:"max_high"`
	MinLow          int64     `json:"min_low"`
	OpeningHigh     int64     `json:"opening_high"`
	OpeningLow      int64     `json:"opening_low"`
	ClosingHigh     int64     `json:"closing_high"`
	ClosingLow      int64     `json:"closing_low"`
	TotalHighVolume int64     `json:"total_high_volume"`
	TotalLowVolume  int64     `json:"total_low_volume"`
	DataPoints      int       `json:"data_points"`
	HourTimestamp   time.Time `gorm:"index:idx_hourly_item_time;not null" json:"hour_timestamp"`
	CreatedAt       time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (PriceHistoryHourly) TableName() string {
	return "price_history_hourly"
}

// PriceHistoryDaily stores daily aggregated price data
type PriceHistoryDaily struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ItemID          int       `gorm:"index:idx_daily_item_date;not null" json:"item_id"`
	AvgHigh         int64     `json:"avg_high"`
	AvgLow          int64     `json:"avg_low"`
	MaxHigh         int64     `json:"max_high"`
	MinLow          int64     `json:"min_low"`
	OpeningHigh     int64     `json:"opening_high"`
	OpeningLow      int64     `json:"opening_low"`
	ClosingHigh     int64     `json:"closing_high"`
	ClosingLow      int64     `json:"closing_low"`
	TotalHighVolume int64     `json:"total_high_volume"`
	TotalLowVolume  int64     `json:"total_low_volume"`
	Volatility      float64   `json:"volatility"`
	DataPoints      int       `json:"data_points"`
	DayDate         time.Time `gorm:"type:date;uniqueIndex:idx_daily_item_date;not null" json:"day_date"`
	CreatedAt       time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (PriceHistoryDaily) TableName() string {
	return "price_history_daily"
}