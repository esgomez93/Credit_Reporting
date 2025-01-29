package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"maas/internal/store"
	"maas/pkg/service"
	mock_service "maas/pkg/service/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetMemesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemeService := mock_service.NewMockMemeService(ctrl)
	memeHandler := NewMemeHandler(mockMemeService)

	t.Run("Successful Request", func(t *testing.T) {
		// Set up expectations for the mock service
		expectedMeme := &store.MemeResponse{
			Meme: "Test meme",
		}
		mockMemeService.EXPECT().
			GetMeme(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedMeme, nil)

		// Create a request
		req := httptest.NewRequest("GET", "/memes?lat=123&lon=456&query=test", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetMemes(w, req)

		// Check the response
		assert.Equal(t, http.StatusOK, w.Code)
		var response store.MemeResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Test meme", response.Meme)
	})

	t.Run("Insufficient Tokens", func(t *testing.T) {
		// Set up expectations for the mock service
		mockMemeService.EXPECT().
			GetMeme(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, service.ErrInsufficientTokens)

		// Create a request
		req := httptest.NewRequest("GET", "/memes?lat=123&lon=456&query=test", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetMemes(w, req)

		// Check the response
		assert.Equal(t, http.StatusPaymentRequired, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Set up expectations for the mock service
		mockMemeService.EXPECT().
			GetMeme(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, service.ErrInvalidAuthToken)

		// Create a request
		req := httptest.NewRequest("GET", "/memes?lat=123&lon=456&query=test", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetMemes(w, req)

		// Check the response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		// Set up expectations for the mock service
		mockMemeService.EXPECT().
			GetMeme(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("some error"))

		// Create a request
		req := httptest.NewRequest("GET", "/memes?lat=123&lon=456&query=test", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetMemes(w, req)

		// Check the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddTokensHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemeService := mock_service.NewMockMemeService(ctrl)
	memeHandler := NewMemeHandler(mockMemeService)

	t.Run("Successful Request", func(t *testing.T) {
		// Set up expectations for the mock service
		mockMemeService.EXPECT().
			AddTokens("test_token", 100).
			Return(nil)

		// Create a request body
		body := map[string]int{"amount": 100}
		jsonBody, _ := json.Marshal(body)

		// Create a request
		req := httptest.NewRequest("POST", "/addtokens", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.AddTokens(w, req)

		// Check the response
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Auth Token", func(t *testing.T) {
		// Create a request body
		body := map[string]int{"amount": 100}
		jsonBody, _ := json.Marshal(body)

		// Create a request without an Authorization header
		req := httptest.NewRequest("POST", "/addtokens", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.AddTokens(w, req)

		// Check the response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		// Create a request with an invalid body
		req := httptest.NewRequest("POST", "/addtokens", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.AddTokens(w, req)

		// Check the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Set up expectations for the mock service to return an error
		mockMemeService.EXPECT().
			AddTokens("test_token", 100).
			Return(errors.New("some error"))

		// Create a request body
		body := map[string]int{"amount": 100}
		jsonBody, _ := json.Marshal(body)

		// Create a request
		req := httptest.NewRequest("POST", "/addtokens", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.AddTokens(w, req)

		// Check the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetBalanceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemeService := mock_service.NewMockMemeService(ctrl)
	memeHandler := NewMemeHandler(mockMemeService)

	t.Run("Successful Request", func(t *testing.T) {
		// Set up expectations for the mock service
		mockMemeService.EXPECT().
			GetTokenBalance("test_token").
			Return(100, nil)

		// Create a request
		req := httptest.NewRequest("GET", "/balance", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetBalance(w, req)

		// Check the response
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]int
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 100, response["token_balance"])
	})

	t.Run("Missing Auth Token", func(t *testing.T) {
		// Create a request without an Authorization header
		req := httptest.NewRequest("GET", "/balance", nil)
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetBalance(w, req)

		// Check the response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Set up expectations for the mock service to return an error
		mockMemeService.EXPECT().
			GetTokenBalance("test_token").
			Return(0, errors.New("some error"))

		// Create a request
		req := httptest.NewRequest("GET", "/balance", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the handler
		memeHandler.GetBalance(w, req)

		// Check the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemeService := mock_service.NewMockMemeService(ctrl)
	memeHandler := NewMemeHandler(mockMemeService)

	t.Run("Successful Authentication", func(t *testing.T) {
		// Set up expectations for the mock service
		mockMemeService.EXPECT().
			CheckTokenBalance("test_token").
			Return(nil)

		// Create a mock next handler that simulates a successful request
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Create a request with an Authorization header
		req := httptest.NewRequest("GET", "/some-protected-route", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Call the middleware
		authMiddleware := memeHandler.AuthMiddleware(nextHandler)
		authMiddleware.ServeHTTP(w, req)

		// Check the response
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Auth Token", func(t *testing.T) {
		// Create a request without an Authorization header
		req := httptest.NewRequest("GET", "/some-protected-route", nil)
		w := httptest.NewRecorder()

		// Create a mock next handler (should not be called)
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		// Call the middleware
		authMiddleware := memeHandler.AuthMiddleware(nextHandler)
		authMiddleware.ServeHTTP(w, req)

		// Check the response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Insufficient Tokens", func(t *testing.T) {
		// Set up expectations for the mock service to return ErrInsufficientTokens
		mockMemeService.EXPECT().
			CheckTokenBalance("test_token").
			Return(service.ErrInsufficientTokens)

		// Create a request with an Authorization header
		req := httptest.NewRequest("GET", "/some-protected-route", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Create a mock next handler (should not be called)
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		// Call the middleware
		authMiddleware := memeHandler.AuthMiddleware(nextHandler)
		authMiddleware.ServeHTTP(w, req)

		// Check the response
		assert.Equal(t, http.StatusPaymentRequired, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Set up expectations for the mock service to return an error
		mockMemeService.EXPECT().
			CheckTokenBalance("test_token").
			Return(errors.New("some error"))

		// Create a request with an Authorization header
		req := httptest.NewRequest("GET", "/some-protected-route", nil)
		req.Header.Set("Authorization", "test_token")
		w := httptest.NewRecorder()

		// Create a mock next handler (should not be called)
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		// Call the middleware
		authMiddleware := memeHandler.AuthMiddleware(nextHandler)
		authMiddleware.ServeHTTP(w, req)

		// Check the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
