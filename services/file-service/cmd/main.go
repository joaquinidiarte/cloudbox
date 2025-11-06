package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/file-service/internal/handler"
	"github.com/joaquinidiarte/cloudbox/services/file-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/services/file-service/internal/service"
	"github.com/joaquinidiarte/cloudbox/shared/config"
	"github.com/joaquinidiarte/cloudbox/shared/middleware"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()
	logger := utils.NewLogger("file-service")

	if err := os.MkdirAll(cfg.StoragePath, 0755); err != nil {
		log.Fatal("Failed to create storage directory", err)
	}

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
	logger.Info("Connected to MongoDB successfully")

	db := client.Database(cfg.MongoDatabase)

	// Init JWT
	tokenDuration, _ := time.ParseDuration(cfg.JWTExpiration)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, tokenDuration)

	// Layers
	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, cfg.StoragePath, cfg.MaxFileSize)
	fileHandler := handler.NewFileHandler(fileService, logger)

	// Init Gin router
	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "file-service"})
	})

	v1 := router.Group("/api/v1/files")
	v1.Use(middleware.AuthMiddleware(jwtManager))
	{
		v1.POST("/upload", fileHandler.UploadFile)
		v1.GET("/", fileHandler.ListFiles)
		v1.GET("/:id/download", fileHandler.DownloadFile)
		v1.DELETE("/:id", fileHandler.DeleteFile)

		// Folder operations
		v1.POST("/folders", fileHandler.CreateFolder)
		v1.GET("/folders/:id", fileHandler.GetFolderContents)

		// Version operations
		v1.GET("/:id/versions", fileHandler.GetFileVersions)
		v1.GET("/:id/versions/:version/download", fileHandler.DownloadFileVersion)
		v1.POST("/:id/versions/:version/restore", fileHandler.RestoreFileVersion)
		v1.DELETE("/:id/versions/:version", fileHandler.DeleteFileVersion)
	}

	port := cfg.ServicePort
	if port == "" {
		port = "8083"
	}
	logger.Infof("Starting file service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
