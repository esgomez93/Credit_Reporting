package repository

import (
	"database/sql"
	"log"

	"maas/internal/store"
)

// MemeRepository handles database operations for memes.
type MemeRepository struct {
	db *sql.DB
}

// NewMemeRepository creates a new MemeRepository.
func NewMemeRepository(db *sql.DB) *MemeRepository {
	return &MemeRepository{
		db: db,
	}
}

// DeductToken decrements a client's token balance by 1.
func (r *MemeRepository) DeductToken(authToken string) error {
	_, err := r.db.Exec("UPDATE clients SET token_balance = token_balance - 1 WHERE auth_token = $1", authToken)
	return err
}

// LogAPICall records an API call in the database.
func (r *MemeRepository) LogAPICall(authToken string) error {
	var clientID int
	err := r.db.QueryRow("SELECT client_id FROM clients WHERE auth_token = $1", authToken).Scan(&clientID)
	if err != nil {
		return err // Handle appropriately, possibly logging the error
	}

	_, err = r.db.Exec("INSERT INTO api_calls (client_id) VALUES ($1)", clientID)
	return err
}

// GetTokenBalance retrieves the token balance for a client.
func (r *MemeRepository) GetTokenBalance(authToken string) (int, error) {
	var tokenBalance int
	err := r.db.QueryRow("SELECT token_balance FROM clients WHERE auth_token = $1", authToken).Scan(&tokenBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, service.ErrInvalidAuthToken
		}
		return 0, err
	}
	return tokenBalance, nil
}

// AddTokens adds tokens to a client's balance.
func (r *MemeRepository) AddTokens(authToken string, amount int) error {
	_, err := r.db.Exec("UPDATE clients SET token_balance = token_balance + $1 WHERE auth_token = $2", amount, authToken)
	return err
}
