package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/auth-service/internal/handler"
	"github.com/joaquinidiarte/cloudbox/services/auth-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/services/auth-service/internal/service"
	"github.com/joaquinidiarte/cloudbox/shared/config"
	"github.com/joaquinidiarte/cloudbox/shared/middleware"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg := config.Load()
	logger := utils.NewLogger("auth-service")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Ping database
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	logger.Info("Connected to MongoDB")

	// Initialize database
	db := mongoClient.Database(cfg.MongoDatabase)

	// Initialize JWT manager
	tokenDuration, _ := time.ParseDuration(cfg.JWTExpiration)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, tokenDuration)

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService, logger)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "auth-service"})
	})

	// Auth routes
	v1 := router.Group("/api/v1/auth")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/verify", authHandler.VerifyToken)
	}

	// Start server
	port := cfg.ServicePort
	if port == "" {
		port = "8081"
	}
	logger.Infof("Auth service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
