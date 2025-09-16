package helper

import (
	"crypto/rand"
	"jti-super-app-go/config"
	"jti-super-app-go/pkg/constants"
	"time"

	"github.com/gin-gonic/gin"
)

func CSRFToken(sid string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := b64url(b)
	return token, nil
}

func ValidateCSRF(c *gin.Context) bool {
	token := c.PostForm("csrf_token")
	if token == "" {
		return false
	}
	sid, err := c.Cookie(constants.CSRF_ID_TOKEN)
	if err != nil || sid == "" {
		return false
	}
	rc, cancel := c.Request.Context(), func() {}
	defer cancel()

	removed, err := config.Rdb.SRem(rc, "csrf:sess:"+sid, token).Result()
	if err != nil {
		return false
	}
	if removed != 1 {
		return false // token sudah dipakai / tidak ada ⇒ invalid
	}
	// Refresh TTL (opsional)
	_ = config.Rdb.Expire(rc, "csrf:sess:"+sid, 30*time.Minute).Err()
	return true
}
