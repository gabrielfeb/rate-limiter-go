package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("should return default values when env vars are not set", func(t *testing.T) {
		// Limpa as variáveis de ambiente
		os.Clearenv()

		cfg := NewConfig()

		assert.Equal(t, "localhost:6379", cfg.RedisAddr)
		assert.Equal(t, 5, cfg.IPLimit)
		assert.Equal(t, time.Minute, cfg.IPBlockDuration)
	})

	t.Run("should return values from environment variables when set", func(t *testing.T) {
		os.Clearenv()

		// Usa t.Setenv para definir variáveis de ambiente temp para o teste
		t.Setenv("REDIS_ADDR", "test-redis:1234")
		t.Setenv("IP_LIMIT_PER_SECOND", "100")
		t.Setenv("TOKEN_BLOCK_DURATION", "10h")

		cfg := NewConfig()

		assert.Equal(t, "test-redis:1234", cfg.RedisAddr)
		assert.Equal(t, 100, cfg.IPLimit)
		assert.Equal(t, 10*time.Hour, cfg.TokenBlockDuration)
	})
}
