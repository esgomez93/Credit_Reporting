package api

import (
	"encoding/json"
	"net/http"

	"maas/pkg/service"
)

// MemeHandler handles API requests related to memes.
type MemeHandler struct {
	memeService *service.MemeService
}

// NewMemeHandler creates a new MemeHandler.
func NewMemeHandler(memeService *service.MemeService) *MemeHandler {
	return &MemeHandler{
		memeService: memeService,
	}
}

// GetMemes handles the GET /memes request.
func (h *MemeHandler) GetMemes(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters (lat, lon, query)
	latitude := r.URL.Query().Get("lat")
	longitude := r.URL.Query().Get("lon")
	query := r.URL.Query().Get("query")

	// Get the auth token from the request header.
	authToken := r.Header.Get("Authorization")

	// Fetch the meme using the service layer.
	meme, err := h.memeService.GetMeme(latitude, longitude, query, authToken)
	if err != nil {
		// Handle errors appropriately (e.g., insufficient tokens, invalid token)
		switch err {
		case service.ErrInsufficientTokens:
			http.Error(w, "Insufficient token balance", http.StatusPaymentRequired)
		case service.ErrInvalidAuthToken:
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the meme data.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meme)
}

// AddTokens handles adding tokens to a client's balance.
func (h *MemeHandler) AddTokens(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	var req service.AddTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.memeService.AddTokens(authToken, req.Amount); err != nil {
		http.Error(w, "Failed to add tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tokens added successfully"))
}

// GetBalance handles retrieving a client's token balance.
func (h *MemeHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	balance, err := h.memeService.GetTokenBalance(authToken)
	if err != nil {
		http.Error(w, "Failed to get token balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"token_balance": balance})
}
