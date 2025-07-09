package limiter

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

// Implementação de Storage para testes unitários
type MockStorage struct {
	mu          sync.Mutex
	counts      map[string]int64
	blocked     map[string]time.Time
	expirations map[string]time.Time

	ShouldErrorOnIncrement bool
	ShouldErrorOnIsBlocked bool
}

// Cria uma nova instância de MockStorage
func NewMockStorage() *MockStorage {
	return &MockStorage{
		counts:      make(map[string]int64),
		blocked:     make(map[string]time.Time),
		expirations: make(map[string]time.Time),
	}
}

// Incrementa o contador para uma chave específica e define a expiração
func (m *MockStorage) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ShouldErrorOnIncrement {
		return 0, errors.New("mock storage increment error")
	}

	if exp, ok := m.expirations[key]; ok && time.Now().After(exp) {
		delete(m.counts, key)
		delete(m.expirations, key)
	}

	m.counts[key]++
	m.expirations[key] = time.Now().Add(window)
	return m.counts[key], nil
}

// Bloqueia uma chave por um período de tempo específico.
func (m *MockStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blocked[key] = time.Now().Add(duration)
	return nil
}

// Verifica se uma chave está bloqueada.
func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ShouldErrorOnIsBlocked {
		return false, errors.New("mock storage isblocked error")
	}

	if exp, ok := m.blocked[key]; ok {
		if time.Now().Before(exp) {
			return true, nil
		}
		// Bloqueio expirou, remove a chave
		delete(m.blocked, key)

		// Reseta a contagem de requisições
		requestKey := strings.Replace(key, "limiter:blocked:", "limiter:requests:", 1)
		delete(m.counts, requestKey)
		delete(m.expirations, requestKey)
	}
	return false, nil
}
