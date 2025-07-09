package middleware

import (
	"log"
	"net"
	"net/http"

	"rate-limiter-go/limiter"
)

func RateLimitMiddleware(limiter *limiter.RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var identifier string
			var isToken bool

			// Determina o identificador da requisição
			token := r.Header.Get("API_KEY")
			if token != "" {
				identifier = token
				isToken = true
			} else {
				// Obtem o IP do cliente a partir do RemoteAddr
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					// Se não conseguir obter o IP, usa o RemoteAddr como identificador
					log.Printf("Could not parse IP from %s: %v", r.RemoteAddr, err)
					identifier = r.RemoteAddr
				} else {
					identifier = ip
				}
				isToken = false
			}

			// Loga o identificador e o tipo de limite
			if !limiter.Allow(r.Context(), identifier, isToken) {
				http.Error(
					w,
					"you have reached the maximum number of requests or actions allowed within a certain time frame",
					http.StatusTooManyRequests,
				)
				return
			}

			// Se a requisição for permitida, continua para o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}
