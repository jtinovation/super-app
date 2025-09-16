package helper

import (
	"crypto/rand"
	"encoding/base64"
	"html"
	"jti-super-app-go/config"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func SanitizeInput(input string) string {
	return strings.TrimSpace(html.EscapeString(input))
}

func GetUrlFile(path, filename string) string {
	if path == "" && filename == "" {
		return ""
	}

	fullPath := path
	if filename != "" {
		if strings.HasSuffix(path, "/") {
			fullPath = path + filename
		} else {
			fullPath = path + "/" + filename
		}
	}

	endpoint := config.AppConfig.Minio.URL
	bucket := config.AppConfig.Minio.Bucket

	return endpoint + "/" + bucket + fullPath
}

func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func GenCode() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return b64url(b)
}

func RedirectBackToLogin(c *gin.Context, loginPath, returnTo, issuer, errMsg string) {
	// if !IsSafeReturnTo(returnTo, issuer) {
	// 	returnTo = "/oauth/authorize" // fallback aman
	// }
	u := url.URL{Path: loginPath}
	q := u.Query()
	q.Set("return_to", returnTo)
	if errMsg != "" {
		q.Set("error", b64url([]byte(errMsg)))
	}
	u.RawQuery = q.Encode()
	c.Redirect(302, u.String())
}
