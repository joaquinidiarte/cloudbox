package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/user-service/internal/handler"
	"github.com/joaquinidiarte/cloudbox/services/user-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/services/user-service/internal/service"
	"github.com/joaquinidiarte/cloudbox/shared/config"
	"github.com/joaquinidiarte/cloudbox/shared/middleware"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg := config.Load()
	logger := utils.NewLogger("user-service")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	logger.Info("Connected to MongoDB")

	db := client.Database(cfg.MongoDatabase)

	// Initialize JWT manager
	tokenDuration, _ := time.ParseDuration(cfg.JWTExpiration)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, tokenDuration)

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "user-service"})
	})

	// User routes (protected)
	v1 := router.Group("/api/v1/users")
	v1.Use(middleware.AuthMiddleware(jwtManager))
	{
		v1.GET("/me", userHandler.GetCurrentUser)
		v1.PUT("/me", userHandler.UpdateCurrentUser)
		v1.POST("/storage", userHandler.UpdateStorageUsed)
		v1.GET("/:id", userHandler.GetUserByID)
	}

	// Start server
	port := cfg.ServicePort
	if port == "" {
		port = "8082"
	}
	logger.Infof("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
