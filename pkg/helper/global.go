package helper

import (
	"html"
	"jti-super-app-go/config"
	"strings"
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
