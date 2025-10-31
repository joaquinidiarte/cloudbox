package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	LogLevel    string
	
	// Service specific
	ServiceName string
	ServicePort string
	GRPCPort    string
	
	// Database
	MongoURI      string
	MongoDatabase string
	
	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	
	// JWT
	JWTSecret     string
	JWTExpiration string
	
	// File Storage
	StoragePath string
	MaxFileSize int64
	
	// API Gateway
	APIGatewayURL string
	
	// Service URLs
	AuthServiceURL string
	UserServiceURL string
	FileServiceURL string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "104857600"), 10, 64)

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		
		ServiceName: getEnv("SERVICE_NAME", "cloudbox-service"),
		ServicePort: getEnv("SERVICE_PORT", "8080"),
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "cloudbox"),
		
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		
		JWTSecret:     getEnv("JWT_SECRET", "change-this-secret-key"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
		
		StoragePath: getEnv("STORAGE_PATH", "./storage"),
		MaxFileSize: maxFileSize,
		
		APIGatewayURL: getEnv("API_GATEWAY_URL", "http://localhost:8080"),
		
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "localhost:50051"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "localhost:50052"),
		FileServiceURL: getEnv("FILE_SERVICE_URL", "localhost:50053"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}