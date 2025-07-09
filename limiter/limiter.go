package limiter

import (
	"context"
	"fmt"
	"time"

	"rate-limiter-go/config"
)

// Define o RateLimiter que gerencia as requisições
type RateLimiter struct {
	storage Storage
	config  *config.Config
}

// Cria uma nova instância de RateLimiter
func NewRateLimiter(storage Storage, cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  cfg,
	}
}

// Verifica se uma requisição é permitida
func (rl *RateLimiter) Allow(ctx context.Context, identifier string, isToken bool) bool {
	limit, blockDuration := rl.getLimits(isToken)

	blockKey := fmt.Sprintf("limiter:blocked:%s", identifier)
	if blocked, err := rl.storage.IsBlocked(ctx, blockKey); err != nil || blocked {
		return false
	}

	requestKey := fmt.Sprintf("limiter:requests:%s", identifier)

	count, err := rl.storage.Increment(ctx, requestKey, 1*time.Second)
	if err != nil {
		return false
	}

	if count > int64(limit) {
		rl.storage.Block(ctx, blockKey, blockDuration)
		return false
	}

	return true
}

// Retorna os limites de requisições e o tempo de bloqueio
func (rl *RateLimiter) getLimits(isToken bool) (int, time.Duration) {
	if isToken {
		return rl.config.TokenLimit, rl.config.TokenBlockDuration
	}
	return rl.config.IPLimit, rl.config.IPBlockDuration
}
