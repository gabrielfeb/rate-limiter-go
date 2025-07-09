package limiter

import (
	"context"
	"testing"
	"time"

	"rate-limiter-go/config"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	cfg := &config.Config{
		IPLimit:            3,
		IPBlockDuration:    60 * time.Millisecond,
		TokenLimit:         5,
		TokenBlockDuration: 120 * time.Millisecond,
	}
	ctx := context.Background()

	testCases := []struct {
		name          string
		identifier    string
		isToken       bool
		setupLimiter  func(t *testing.T) *RateLimiter
		preAssert     func(t *testing.T, l *RateLimiter)
		expectedAllow bool
	}{
		// Testes de IP
		{
			name:       "IP under limit",
			identifier: "192.168.1.1",
			isToken:    false,
			setupLimiter: func(t *testing.T) *RateLimiter {
				return NewRateLimiter(NewMockStorage(), cfg)
			},
			preAssert: func(t *testing.T, l *RateLimiter) {
				l.Allow(ctx, "192.168.1.1", false) // 1
				l.Allow(ctx, "192.168.1.1", false) // 2
			},
			expectedAllow: true, // Testa a 3ª requisição
		},
		{
			name:       "IP exceeds limit",
			identifier: "192.168.1.2",
			isToken:    false,
			setupLimiter: func(t *testing.T) *RateLimiter {
				return NewRateLimiter(NewMockStorage(), cfg)
			},
			preAssert: func(t *testing.T, l *RateLimiter) {
				for i := 0; i < cfg.IPLimit; i++ {
					l.Allow(ctx, "192.168.1.2", false)
				}
			},
			expectedAllow: false, // Testa a 4ª requisição
		},
		// Testes de Token
		{
			name:       "Token under limit",
			identifier: "token-abc",
			isToken:    true,
			setupLimiter: func(t *testing.T) *RateLimiter {
				return NewRateLimiter(NewMockStorage(), cfg)
			},
			preAssert: func(t *testing.T, l *RateLimiter) {
				for i := 0; i < cfg.TokenLimit-1; i++ {
					l.Allow(ctx, "token-abc", true)
				}
			},
			expectedAllow: true, // Testa a 5ª requisição
		},
		{
			name:       "Token exceeds limit",
			identifier: "token-def",
			isToken:    true,
			setupLimiter: func(t *testing.T) *RateLimiter {
				return NewRateLimiter(NewMockStorage(), cfg)
			},
			preAssert: func(t *testing.T, l *RateLimiter) {
				for i := 0; i < cfg.TokenLimit; i++ {
					l.Allow(ctx, "token-def", true)
				}
			},
			expectedAllow: false, // Testa a 6ª requisição
		},
		// Testes de Bloqueio
		{
			name:       "Identifier is already blocked",
			identifier: "192.168.1.3",
			isToken:    false,
			setupLimiter: func(t *testing.T) *RateLimiter {
				mock := NewMockStorage()
				// Bloqueia manualmente o identificador antes do teste
				mock.Block(ctx, "limiter:blocked:192.168.1.3", cfg.IPBlockDuration)
				return NewRateLimiter(mock, cfg)
			},
			expectedAllow: false,
		},
		{
			name:       "Block expires and allows new request",
			identifier: "192.168.1.4",
			isToken:    false,
			setupLimiter: func(t *testing.T) *RateLimiter {
				limiter := NewRateLimiter(NewMockStorage(), cfg)
				// Excede o limite para causar um bloqueio
				for i := 0; i < cfg.IPLimit+1; i++ {
					limiter.Allow(ctx, "192.168.1.4", false)
				}
				assert.False(t, limiter.Allow(ctx, "192.168.1.4", false), "deveria estar bloqueado imediatamente")
				return limiter
			},
			preAssert: func(t *testing.T, l *RateLimiter) {
				time.Sleep(cfg.IPBlockDuration + 10*time.Millisecond)
			},
			expectedAllow: true,
		},
		// Testes de Erro na Storage
		{
			name:       "Storage error on IsBlocked",
			identifier: "192.168.1.5",
			isToken:    false,
			setupLimiter: func(t *testing.T) *RateLimiter {
				mock := NewMockStorage()
				mock.ShouldErrorOnIsBlocked = true // Simula erro
				return NewRateLimiter(mock, cfg)
			},
			expectedAllow: false, // Deve negar a requisição em caso de erro
		},
		{
			name:       "Storage error on Increment",
			identifier: "192.168.1.6",
			isToken:    false,
			setupLimiter: func(t *testing.T) *RateLimiter {
				mock := NewMockStorage()
				mock.ShouldErrorOnIncrement = true // Simula erro
				return NewRateLimiter(mock, cfg)
			},
			expectedAllow: false, // Deve negar a requisição em caso de erro
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			limiter := tc.setupLimiter(t)

			if tc.preAssert != nil {
				tc.preAssert(t, limiter)
			}

			actualAllow := limiter.Allow(ctx, tc.identifier, tc.isToken)
			assert.Equal(t, tc.expectedAllow, actualAllow)
		})
	}
}
