package handler

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/shared/config"
	"github.com/joaquinidiarte/cloudbox/shared/models"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
)

type ProxyHandler struct {
	config *config.Config
	logger *utils.Logger
}

func NewProxyHandler(cfg *config.Config, logger *utils.Logger) *ProxyHandler {
	return &ProxyHandler{
		config: cfg,
		logger: logger,
	}
}

func (h *ProxyHandler) ProxyToAuth(c *gin.Context) {
	// Use service name in Docker, localhost for local development
	baseURL := getEnv("AUTH_SERVICE_URL", "http://auth-service:8081")
	h.proxyRequest(c, baseURL)
}

func (h *ProxyHandler) proxyRequest(c *gin.Context, targetURL string) {
	// Build target URL
	url := targetURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		url += "?" + c.Request.URL.RawQuery
	}

	h.logger.Infof("Proxying %s %s to %s", c.Request.Method, c.Request.URL.Path, url)

	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Create new request
	req, err := http.NewRequest(c.Request.Method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to proxy request"))
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to send request: %v", err)
		c.JSON(http.StatusBadGateway, models.ErrorResponse("Service unavailable"))
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorf("Failed to read response: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to read response"))
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Send response
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
