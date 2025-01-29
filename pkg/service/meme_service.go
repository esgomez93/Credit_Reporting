package service

import (
	"errors"

	"maas/internal/store"
	"maas/pkg/repository"
	"maas/utils"
)

// MemeService handles the business logic for memes.
type MemeService struct {
	memeRepo *repository.MemeRepository
}

// NewMemeService creates a new MemeService.
func NewMemeService(memeRepo *repository.MemeRepository) *MemeService {
	return &MemeService{
		memeRepo: memeRepo,
	}
}

// ErrInsufficientTokens is returned when a client has insufficient tokens.
var ErrInsufficientTokens = errors.New("insufficient tokens")

// ErrInvalidAuthToken is returned when the provided auth token is invalid.
var ErrInvalidAuthToken = errors.New("invalid authorization token")

// GetMeme fetches a meme, checks token balance, and updates it.
func (s *MemeService) GetMeme(latitude, longitude, query, authToken string) (*store.MemeResponse, error) {
	// Check if the client has enough tokens.
	if err := s.CheckTokenBalance(authToken); err != nil {
		return nil, err
	}

	// Deduct token for the API call.
	if err := s.memeRepo.DeductToken(authToken); err != nil {
		return nil, err
	}

	// Log the API call.
	if err := s.memeRepo.LogAPICall(authToken); err != nil {
		// Log the error, but don't fail the request.
	}

	// Generate a random meme using the utility function.
	memeText := utils.GenerateRandomMeme(query)

	// Create a MemeResponse object.
	meme := &store.MemeResponse{
		Meme:      memeText,
		Latitude:  latitude,
		Longitude: longitude,
		Query:     query,
	}

	return meme, nil
}

// CheckTokenBalance checks if the client has a sufficient token balance.
func (s *MemeService) CheckTokenBalance(authToken string) error {
	tokenBalance, err := s.memeRepo.GetTokenBalance(authToken)
	if err != nil {
		return err
	}

	if tokenBalance <= 0 {
		return ErrInsufficientTokens
	}

	return nil
}

// AddTokensRequest represents the request body for adding tokens.
type AddTokensRequest struct {
	Amount int `json:"amount"`
}

// AddTokens adds tokens to a client's balance.
func (s *MemeService) AddTokens(authToken string, amount int) error {
	return s.memeRepo.AddTokens(authToken, amount)
}

// GetTokenBalance retrieves the token balance for a client.
func (s *MemeService) GetTokenBalance(authToken string) (int, error) {
	return s.memeRepo.GetTokenBalance(authToken)
}
