package store

import "time"

// Client represents a client in the database.
type Client struct {
	ClientID     int    `db:"client_id"`
	AuthToken    string `db:"auth_token"`
	TokenBalance int    `db:"token_balance"`
}

// APICall represents an API call made by a client.
type APICall struct {
	CallID    int       `db:"call_id"`
	ClientID  int       `db:"client_id"`
	Timestamp time.Time `db:"timestamp"`
}

// MemeResponse represents the API response structure.
type MemeResponse struct {
	Meme      string `json:"meme"`
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	Query     string `json:"query,omitempty"`
}
