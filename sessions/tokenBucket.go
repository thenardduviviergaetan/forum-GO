package forum

import (
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	tokens              int
	maxTokens           int
	refillRate          time.Duration
	lastRefillTimestamp time.Time
	mutex               sync.Mutex
}

// Initialize a new token bucket
func NewBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:              maxTokens,
		maxTokens:           maxTokens,
		refillRate:          refillRate,
		lastRefillTimestamp: time.Now(),
	}
}

// Check if request is Allowed. If true, sub token to bucket
func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	elapsed := time.Now().Sub(tb.lastRefillTimestamp)
	tokenToAdd := int(elapsed / tb.refillRate)

	if tokenToAdd > 0 {
		if tb.tokens+tokenToAdd < tb.maxTokens {
			tb.tokens = tb.tokens + tokenToAdd
		} else {
			tb.tokens = tb.maxTokens
		}
	}
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// Custom handler with token bucket limiter for rate limiting
func HandleWithLimiter(path string, handler func(http.ResponseWriter, *http.Request), limiter *TokenBucket) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		handler(w, r)
	})
}
