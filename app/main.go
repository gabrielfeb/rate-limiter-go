package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"rate-limiter-go/config"
	"rate-limiter-go/limiter"
	ratelimit "rate-limiter-go/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	cfg := config.NewConfig()
	log.Println("Configuration loaded.")

	// Inicialização do Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis.")

	redisStorage := limiter.NewRedisStorage(rdb)
	rateLimiter := limiter.NewRateLimiter(redisStorage, cfg)
	log.Println("Rate limiter initialized.")

	// Configura Chi Router
	r := chi.NewRouter()

	// Configura Chi Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Configura nosso Rate Limiter Middleware
	log.Println("Applying rate limiter middleware...")
	r.Use(ratelimit.RateLimitMiddleware(rateLimiter))

	// Definição de rotas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Welcome! Your request was successful."}`)
	})

	// Rota exemplo
	r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"user": "Alex", "profile": "expert"}`)
	})

	// Inicialização do Server
	log.Println("Starting server on port 8080 with Chi router...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
