package models

import "time"

// Stock represents a stock analyst rating/recommendation from the API.
type Stock struct {
	ID         int       `json:"id"`
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`
	Brokerage  string    `json:"brokerage"`
	Action     string    `json:"action"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	TargetFrom float64   `json:"target_from"`
	TargetTo   float64   `json:"target_to"`
	CreatedAt  time.Time `json:"created_at"`
}

// StockListResponse represents the paginated response from our API.
type StockListResponse struct {
	Items      []Stock `json:"items"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}

// Recommendation represents a stock recommendation with a score.
type Recommendation struct {
	Stock  Stock   `json:"stock"`
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

// APIResponse represents the response from the KarenAI API.
type APIResponse struct {
	Items    []APIStock `json:"items"`
	NextPage string     `json:"next_page"`
}

// APIStock represents a single stock item from the KarenAI API.
type APIStock struct {
	Ticker     string  `json:"ticker"`
	Company    string  `json:"company"`
	Brokerage  string  `json:"brokerage"`
	Action     string  `json:"action"`
	RatingFrom string  `json:"rating_from"`
	RatingTo   string  `json:"rating_to"`
	TargetFrom float64 `json:"target_from"`
	TargetTo   float64 `json:"target_to"`
}

// LoginRequest represents the login payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response.
type LoginResponse struct {
	AuthToken string `json:"auth_token"`
	Message   string `json:"message"`
}

// SyncStatus represents the status of a data sync operation.
type SyncStatus struct {
	Message      string `json:"message"`
	ItemsSynced  int    `json:"items_synced"`
	PagesScanned int    `json:"pages_scanned"`
}
