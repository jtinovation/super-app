package helper

import (
	"encoding/base64"
	"encoding/json"
	"jti-super-app-go/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

const CookieName = "sso_session"

func SetSSO(c *gin.Context, sub *dto.LoginResponseDTO, maxAgeSeconds int) {
	b, _ := json.Marshal(sub)
	val := base64.RawURLEncoding.EncodeToString(b)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		CookieName,
		val,
		maxAgeSeconds,
		"/",
		// config.AppConfig.CookieDomain, // contoh: .example.com atau sso.example.com
		"",
		true, // secure: wajib true di HTTPS
		true, // httpOnly
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
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(CookieName, "", -1, "/", "", true, true)
}
