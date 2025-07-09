package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Interface que define os métodos necessários para o armazenamento de dados do rate limiter
type Storage interface {
	Increment(ctx context.Context, key string, window time.Duration) (int64, error)
	Block(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}

// Implementação de Storage usando Redis como backend
type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

// Incrementa o contador para uma chave específica e define a expiração
func (r *RedisStorage) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := r.client.TxPipeline()

	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return incr.Val(), nil
}

// Bloqueia uma chave por um período de tempo específico
func (r *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key, "blocked", duration).Err()
}

// Verifica se uma chave está bloqueada
func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return true, err
	}
	return val == 1, nil
}
