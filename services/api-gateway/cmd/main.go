package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/api-gateway/internal/handler"
	"github.com/joaquinidiarte/cloudbox/shared/config"
	"github.com/joaquinidiarte/cloudbox/shared/middleware"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
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

		// File routes
		files := api.Group("/files")
		{
			files.POST("/upload", proxyHandler.ProxyToFile)
			files.GET("/", proxyHandler.ProxyToFile)
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
