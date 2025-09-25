package middleware

import (
	"context"
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rateLimit       = 10          // requests
	rateLimitWindow = time.Minute // per minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Terapkan hanya untuk endpoint login POST
		if c.Request.Method != http.MethodPost || c.FullPath() != "/api/v1/auth/login" {
			c.Next()
			return
		}

		ip := c.ClientIP()

		// Window dalam detik sebagai int64 (hindari campur float64)
		windowSec := int64(rateLimitWindow / time.Second) // mis. 60
		nowSec := time.Now().Unix()
		slot := nowSec / windowSec

		key := fmt.Sprintf("rl:login:%s:%d", ip, slot)

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()

		n, _ := config.Rdb.Incr(ctx, key).Result()
		// TTL sedikit > window untuk safety
		_ = config.Rdb.Expire(ctx, key, rateLimitWindow+2*time.Second).Err()

		if int(n) > rateLimit {
			// Hitung sisa detik hingga window berikutnya
			elapsed := nowSec % windowSec
			reset := int(windowSec - elapsed)
			if reset <= 0 {
				reset = int(windowSec)
			}
			c.Header("Retry-After", strconv.Itoa(reset))
			helper.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", fmt.Errorf("too many requests"))
			c.Abort()
			return
		}

		c.Next()
	}
}
