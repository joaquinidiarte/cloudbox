package main

import (
	"log"

	"github.com/cloudbox/services/api-gateway/internal/handler"
	"github.com/cloudbox/shared/config"
	"github.com/cloudbox/shared/middleware"
	"github.com/cloudbox/shared/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	logger := utils.NewLogger("api-gateway")

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	// Initialize proxy handler
	proxyHandler := handler.NewProxyHandler(cfg, logger)

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", proxyHandler.ProxyToAuth)
			auth.POST("/login", proxyHandler.ProxyToAuth)
			auth.POST("/verify", proxyHandler.ProxyToAuth)
		}
	}

	// Start server
	port := cfg.ServicePort
	if port == "" {
		port = "8080"
	}
	logger.Infof("API Gateway starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}