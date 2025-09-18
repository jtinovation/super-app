package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/internal/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const CookieName = "sso_session"

func SetSSO(c *gin.Context, sub *dto.LoginResponseDTO, maxAgeSeconds int) {
	b, _ := json.Marshal(sub)

	// store b to redis with key sub.User.ID and expiry maxAgeSeconds
	err := config.Rdb.Set(c.Request.Context(), "sso:"+sub.User.ID, b, time.Duration(maxAgeSeconds)*time.Second).Err()
	if err != nil {
		// log error, tapi tidak perlu di-handle lebih lanjut
		fmt.Println("Failed to store SSO session in Redis:", err)
	}

	val := base64.RawURLEncoding.EncodeToString([]byte(sub.User.ID))

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		CookieName,
		val,
		maxAgeSeconds,
		"/",
		config.AppConfig.CookieDomain, // contoh: .example.com atau sso.example.com
		true,                          // secure: wajib true di HTTPS
		true,                          // httpOnly
	)
}

func GetSSO(c *gin.Context) (string, bool) {
	val, err := c.Cookie(CookieName)
	if err != nil || val == "" {
		return "", false
	}
	b, err := base64.RawURLEncoding.DecodeString(val)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func ClearSSO(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(CookieName, "", -1, "/", config.AppConfig.CookieDomain, true, true)
}
