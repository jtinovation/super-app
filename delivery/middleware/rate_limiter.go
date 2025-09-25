package middleware

import (
	"context"
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rateLimit     = 10          // requests
	rateLimitTime = time.Minute // per minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		slot := time.Now().Unix() / int64(rateLimitTime.Seconds())
		key := fmt.Sprintf("rl:login:%s:%d", ip, slot)

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()
		config.Rdb.Incr(ctx, key).Result()

		n, _ := config.Rdb.Incr(ctx, key).Result()
		_ = config.Rdb.Expire(ctx, key, rateLimitTime).Err()

		if int(n) > rateLimit {
			helper.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", fmt.Errorf("too many requests"))
			return
		}
		c.Next()
	}
}
