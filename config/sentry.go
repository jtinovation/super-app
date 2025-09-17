package config

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func InitSentry() {
	if AppConfig.SentryDSN == "" {
		return
	}

	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: AppConfig.SentryDSN,
		// Debug:            true,
		EnableTracing:    true,
		SendDefaultPII:   true,
		TracesSampleRate: 0.8,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	log.Println("Sentry initialized")
}
