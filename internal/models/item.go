package models

// ItemPrice represents the price information for an OSRS item
type ItemPrice struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	High       int64  `json:"high"`
	HighTime   int64  `json:"highTime"`
	Low        int64  `json:"low"`
	LowTime    int64  `json:"lowTime"`
	HighVolume int64  `json:"highVolume,omitempty"` // Volume of high price trades
	LowVolume  int64  `json:"lowVolume,omitempty"`  // Volume of low price trades
}

// OSRSWikiResponse represents the response from the OSRS Wiki API
type OSRSWikiResponse struct {
	Data map[string]ItemPriceData `json:"data"`
}

// ItemPriceData represents the raw price data from the API
type ItemPriceData struct {
	High       int64 `json:"high"`
	HighTime   int64 `json:"highTime"`
	Low        int64 `json:"low"`
	LowTime    int64 `json:"lowTime"`
	HighVolume int64 `json:"highPriceVolume"` // Volume from OSRS Wiki API
	LowVolume  int64 `json:"lowPriceVolume"`  // Volume from OSRS Wiki API
}

// ItemMapping represents item ID to name mapping
type ItemMapping struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}