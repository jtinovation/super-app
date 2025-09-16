package middleware

import (
	"jti-super-app-go/config"
	"jti-super-app-go/pkg/constants"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CSRFTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(constants.CSRF_ID_TOKEN)
		if err != nil || sid == "" {
			sid = uuid.NewString()
			c.SetSameSite(http.SameSiteLaxMode) // Lax cukup untuk form POST dari domain yang sama
			c.SetCookie(
				constants.CSRF_ID_TOKEN,
				sid,
				int(30*time.Minute.Seconds()), // boleh diset lebih panjang; kita perpanjang di Redis juga
				"/",
				"",   // atau domain IdP kamu; jaga konsistensi dengan domain login
				true, // Secure (wajib HTTPS)
				true, // HttpOnly
			)
		}

		token, err := helper.CSRFToken(sid)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate CSRF token", err)
			c.Abort()
			return
		}

		// SADD token ke set session di Redis
		// Simpan token di Redis dengan TTL sesuai cookie
		err = config.Rdb.SAdd(c.Request.Context(), "csrf:sess:"+sid, token).Err()
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to store CSRF token", err)
			c.Abort()
			return
		}
		err = config.Rdb.Expire(c.Request.Context(), "csrf:sess:"+sid, 30*time.Minute).Err()
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to set CSRF token expiry", err)
			c.Abort()
			return
		}

		c.Set("csrf_token", token)
		c.Next()
	}
}
