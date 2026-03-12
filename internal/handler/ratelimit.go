package handler

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter implements a per-IP rate limiter.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) setHeaders(w http.ResponseWriter, limiter *rate.Limiter) {
	tokens := limiter.Tokens()
	remaining := int(math.Max(0, math.Floor(tokens)))
	deficit := float64(rl.burst) - tokens
	var resetUnix int64
	if deficit <= 0 {
		resetUnix = time.Now().Unix()
	} else {
		secsToFull := deficit / float64(rl.rate)
		resetUnix = time.Now().Add(time.Duration(secsToFull * float64(time.Second))).Unix()
	}
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetUnix))
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
		limiter := rl.getLimiter(ip)
		if !limiter.Allow() {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second/time.Duration(rl.rate)).Unix()))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"code":"RATE_LIMITED","message":"Too many requests, please try again later"}}`))
			return
		}
		rl.setHeaders(w, limiter)
		next.ServeHTTP(w, r)
	})
}
