package limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	tokens   int
	buckets  map[int64]*tokenBucket // userID -> bucket
	mu       sync.Mutex
	interval time.Duration
}

type tokenBucket struct {
	tokens    int
	lastReset time.Time
}

func NewRateLimiter(tokens int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:   tokens,
		buckets:  make(map[int64]*tokenBucket),
		interval: interval,
	}
}

func (rl *RateLimiter) Allow(userID int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[userID]
	now := time.Now()

	if !exists || now.Sub(bucket.lastReset) > rl.interval {
		// Новый bucket: 10 созданий в минуту
		rl.buckets[userID] = &tokenBucket{
			tokens:    rl.tokens,
			lastReset: now,
		}
		return true
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}
