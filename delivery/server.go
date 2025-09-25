package delivery

import (
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/delivery/middleware"
	"jti-super-app-go/internal/service"
	"log"

	sentrygin "github.com/getsentry/sentry-go/gin"

	"github.com/gin-gonic/gin"
)

type AppServer struct {
	router *gin.Engine
}

func Server() *AppServer {
	config.LoadConfig()
	config.ConnectDatabase()
	config.ConnectRedis()
	config.InitMinio()
	config.InitSentry()
	db := config.DB

	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*")
	jwtService := service.NewJWTService()

	container := InitContainer(db, jwtService)
	middleware.CORS(router)
	SetupRoutes(router, container, jwtService)

	return &AppServer{router: router}
}

func (s *AppServer) Run() {
	cfg := config.AppConfig
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddr)
	if err := s.router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
