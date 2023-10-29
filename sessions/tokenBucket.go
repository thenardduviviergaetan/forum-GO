package forum

import (
	"net/http"
	"time"
)

type TokenBucket struct {
	tokens              chan struct{}
	maxTokens           int
	refillRate          time.Duration
	lastRefillTimestamp time.Time
}

// Initialize a new token bucket
func NewBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		tokens:     make(chan struct{}, maxTokens),
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}

	for i := 0; i < maxTokens; i++ {
		tb.tokens <- struct{}{}
	}
	go tb.refill()
	return tb
}

// Refill bucket
func (tb *TokenBucket) refill() {
	ticker := time.NewTicker(tb.refillRate)
	for range ticker.C {
		select {
		case tb.tokens <- struct{}{}:
		default:
		}
	}
}

// Check if request is Allowed. If true, sub token to bucket
func (tb *TokenBucket) allow() bool {
	select {
	case <-tb.tokens:
		return true
	default:
		return false
	}
}

// Custom handler with token bucket limiter for rate limiting
func HandleWithLimiter(path string, handler func(http.ResponseWriter, *http.Request), limiter *TokenBucket) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if !limiter.allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		handler(w, r)
	})
}

// type TokenBucket struct {
// 	tokens              int
// 	maxTokens           int
// 	refillRate          time.Duration
// 	lastRefillTimestamp time.Time
// 	mutex               sync.Mutex
// }

// // Initialize a new token bucket
// func NewBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
// 	return &TokenBucket{
// 		tokens:              maxTokens,
// 		maxTokens:           maxTokens,
// 		refillRate:          refillRate,
// 		lastRefillTimestamp: time.Now(),
// 	}
// }

// // Check if request is Allowed. If true, sub token to bucket
// func (tb *TokenBucket) Allow() bool {
// 	tb.mutex.Lock()
// 	defer tb.mutex.Unlock()

// 	elapsed := time.Now().Sub(tb.lastRefillTimestamp)
// 	tokenToAdd := int(elapsed / tb.refillRate)

// 	if tokenToAdd > 0 {
// 		if tb.tokens+tokenToAdd < tb.maxTokens {
// 			tb.tokens = tb.tokens + tokenToAdd
// 		} else {
// 			tb.tokens = tb.maxTokens
// 		}
// 	}
// 	if tb.tokens > 0 {
// 		tb.tokens--
// 		return true
// 	}
// 	return false
// }

// // Custom handler with token bucket limiter for rate limiting
// func HandleWithLimiter(path string, handler func(http.ResponseWriter, *http.Request), limiter *TokenBucket) {
// 	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
// 		if !limiter.Allow() {
// 			http.Error(w, "Too many requests", http.StatusTooManyRequests)
// 			return
// 		}
// 		handler(w, r)
// 	})
// }
