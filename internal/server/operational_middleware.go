package server

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type rateLimitEntry struct {
	started time.Time
	count   int
}

type rateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]rateLimitEntry
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, entries: make(map[string]rateLimitEntry)}
}

func (r *rateLimiter) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		now := time.Now()

		r.mu.Lock()
		entry := r.entries[key]
		if entry.started.IsZero() || now.Sub(entry.started) >= r.window {
			entry = rateLimitEntry{started: now}
		}
		entry.count++
		r.entries[key] = entry
		exceeded := entry.count > r.limit
		r.mu.Unlock()

		if exceeded {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(429, gin.H{"success": false, "message": "too many requests"})
			return
		}
		c.Next()
	}
}

func (s *Server) requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func (s *Server) requestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		s.logger.Info().Str("request_id", c.GetString("request_id")).Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).Int("status", c.Writer.Status()).
			Dur("latency", time.Since(started)).Msg("http request")
	}
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
