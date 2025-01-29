package api

import (
	"net/http"

	"maas/pkg/service"
)

// AuthMiddleware checks for a valid auth token and sufficient token balance.
func (h *MemeHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		err := h.memeService.CheckTokenBalance(authToken)
		if err != nil {
			if err == service.ErrInsufficientTokens {
				http.Error(w, "Insufficient token balance", http.StatusPaymentRequired)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
