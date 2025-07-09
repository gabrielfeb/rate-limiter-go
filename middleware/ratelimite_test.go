package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"rate-limiter-go/config"
	"rate-limiter-go/limiter"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	cfg := &config.Config{
		IPLimit:            1,
		TokenLimit:         2,
		IPBlockDuration:    time.Minute,
		TokenBlockDuration: time.Minute,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	t.Run("should allow request under IP limit", func(t *testing.T) {
		mockStorage := limiter.NewMockStorage()
		rateLimiter := limiter.NewRateLimiter(mockStorage, cfg)
		middleware := RateLimitMiddleware(rateLimiter)(nextHandler)

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "OK", rr.Body.String())
	})

	t.Run("should block request over IP limit", func(t *testing.T) {
		mockStorage := limiter.NewMockStorage()
		rateLimiter := limiter.NewRateLimiter(mockStorage, cfg)
		middleware := RateLimitMiddleware(rateLimiter)(nextHandler)

		// Primeira requisição (permitida)
		req1 := httptest.NewRequest("GET", "/", nil)
		req1.RemoteAddr = "192.168.1.2:12345"
		middleware.ServeHTTP(httptest.NewRecorder(), req1)

		// Segunda requisição (deve ser bloqueada)
		req2 := httptest.NewRequest("GET", "/", nil)
		req2.RemoteAddr = "192.168.1.2:12345"
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req2)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Contains(t, rr.Body.String(), "you have reached the maximum number of requests")
	})

	t.Run("should use token limit when API_KEY header is present", func(t *testing.T) {
		mockStorage := limiter.NewMockStorage()
		rateLimiter := limiter.NewRateLimiter(mockStorage, cfg)
		middleware := RateLimitMiddleware(rateLimiter)(nextHandler)

		// 1ª (IP) - OK
		req1 := httptest.NewRequest("GET", "/", nil)
		req1.RemoteAddr = "192.168.1.3:12345"
		middleware.ServeHTTP(httptest.NewRecorder(), req1)

		// 2ª (IP) - Bloqueado
		req2 := httptest.NewRequest("GET", "/", nil)
		req2.RemoteAddr = "192.168.1.3:12345"
		rrIP := httptest.NewRecorder()
		middleware.ServeHTTP(rrIP, req2)
		assert.Equal(t, http.StatusTooManyRequests, rrIP.Code)

		// 3ª (Token) - OK pois o limite do token (2) é maior que o do IP (1)
		req3 := httptest.NewRequest("GET", "/", nil)
		req3.RemoteAddr = "192.168.1.3:12345"
		req3.Header.Set("API_KEY", "my-secret-token")
		rrToken := httptest.NewRecorder()
		middleware.ServeHTTP(rrToken, req3)
		assert.Equal(t, http.StatusOK, rrToken.Code)
	})
}
