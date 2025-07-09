package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Ações de configuração do rate limiter
type Config struct {
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	IPLimit            int
	IPBlockDuration    time.Duration
	TokenLimit         int
	TokenBlockDuration time.Duration
	DefaultTestToken   string
}

// Carrega as configurações do rate limiter a partir de variáveis de ambiente
func NewConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		IPLimit:            getEnvAsInt("IP_LIMIT_PER_SECOND", 5),
		IPBlockDuration:    getEnvAsDuration("IP_BLOCK_DURATION", time.Minute),
		TokenLimit:         getEnvAsInt("TOKEN_LIMIT_PER_SECOND", 20),
		TokenBlockDuration: getEnvAsDuration("TOKEN_BLOCK_DURATION", 5*time.Minute),
		DefaultTestToken:   getEnv("DEFAULT_TEST_TOKEN", "my-secret-token-123"),
	}
}

// Recupera o valor de uma variável de ambiente ou retorna um valor padrão
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Recupera o valor de uma variável de ambiente como inteiro ou retorna um valor padrão
func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

// Recupera o valor de uma variável de ambiente como time.Duration ou retorna um valor padrão
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	strValue := getEnv(key, "")
	if value, err := time.ParseDuration(strValue); err == nil {
		return value
	}
	return fallback
}
