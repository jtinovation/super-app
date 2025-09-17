package main

import (
	"jti-super-app-go/delivery"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	defer sentry.Flush(2 * time.Second)
	delivery.Server().Run()
}
