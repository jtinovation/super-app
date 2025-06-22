package delivery

import (
	"fmt"
	"jti-super-app-go/config"
	"jti-super-app-go/delivery/middleware"
	"jti-super-app-go/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

type AppServer struct {
	router *gin.Engine
}

func Server() *AppServer {
	config.LoadConfig()
	config.ConnectDatabase()
	config.ConnectRedis()
	db := config.DB
	router := gin.Default()
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
