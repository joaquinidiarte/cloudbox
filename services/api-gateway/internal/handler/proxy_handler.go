package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

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
	baseURL := h.config.AuthServiceURL
	h.proxyRequest(c, baseURL)
}

func (h *ProxyHandler) ProxyToUser(c *gin.Context) {
	// Use service name in Docker, localhost for local development
	baseURL := h.config.UserServiceURL
	h.proxyRequest(c, baseURL)
}

func (h *ProxyHandler) ProxyToFile(c *gin.Context) {
	// Use service name in Docker, localhost for local development
	baseURL := h.config.FileServiceURL
	h.proxyRequest(c, baseURL)
}

func (h *ProxyHandler) DeleteFile(c *gin.Context) {
	// Special handler for file deletion that orchestrates storage update
	fileServiceURL := h.config.FileServiceURL
	userServiceURL := h.config.UserServiceURL

	// Build target URL
	url := fileServiceURL + c.Request.URL.Path
	h.logger.Infof("Deleting file at %s", url)

	// Create request to file service
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to proxy request"))
		return
	}

	// Copy Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Send request to file service
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

	// If deletion was successful, update user storage (decrement)
	if resp.StatusCode == http.StatusOK {
		var response struct {
			Success bool `json:"success"`
			Data    struct {
				DeletedSize int64 `json:"deleted_size"`
			} `json:"data"`
		}
		if err := json.Unmarshal(respBody, &response); err == nil && response.Success {
			go h.updateUserStorage(userServiceURL, -response.Data.DeletedSize, c.GetHeader("Authorization"))
		}
	}

	// Send response to client
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (h *ProxyHandler) DeleteFileVersion(c *gin.Context) {
	// Special handler for file version deletion that orchestrates storage update
	fileServiceURL := h.config.FileServiceURL
	userServiceURL := h.config.UserServiceURL

	// Build target URL
	url := fileServiceURL + c.Request.URL.Path
	h.logger.Infof("Deleting file version at %s", url)

	// Create request to file service
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to proxy request"))
		return
	}

	// Copy Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Send request to file service
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

	// If deletion was successful, update user storage (decrement)
	if resp.StatusCode == http.StatusOK {
		var response struct {
			Success bool `json:"success"`
			Data    struct {
				DeletedSize int64 `json:"deleted_size"`
			} `json:"data"`
		}
		if err := json.Unmarshal(respBody, &response); err == nil && response.Success {
			go h.updateUserStorage(userServiceURL, -response.Data.DeletedSize, c.GetHeader("Authorization"))
		}
	}

	// Send response to client
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (h *ProxyHandler) UploadFile(c *gin.Context) {
	// Special handler for file upload that orchestrates storage update
	fileServiceURL := h.config.FileServiceURL
	userServiceURL := h.config.UserServiceURL

	// First, proxy the upload to file service
	url := fileServiceURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		url += "?" + c.Request.URL.RawQuery
	}

	h.logger.Infof("Uploading file to %s", url)

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("No file provided"))
		return
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		h.logger.Errorf("Failed to create form file: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to process file"))
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		h.logger.Errorf("Failed to open file: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to process file"))
		return
	}
	defer fileContent.Close()

	_, err = io.Copy(part, fileContent)
	if err != nil {
		h.logger.Errorf("Failed to copy file: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to process file"))
		return
	}

	// Add parent_id if present
	if parentID := c.PostForm("parent_id"); parentID != "" {
		writer.WriteField("parent_id", parentID)
	}

	writer.Close()

	// Create request to file service
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to proxy request"))
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Copy Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Send request to file service
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

	// If upload was successful, update user storage
	if resp.StatusCode == http.StatusCreated {
		go h.updateUserStorage(userServiceURL, file.Size, c.GetHeader("Authorization"))
	}

	// Send response to client
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (h *ProxyHandler) updateUserStorage(userServiceURL string, fileSize int64, authHeader string) {
	url := userServiceURL + "/api/v1/users/storage"

	payload := map[string]int64{"increment": fileSize}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		h.logger.Errorf("Failed to marshal storage update: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		h.logger.Errorf("Failed to create storage update request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to update user storage: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.logger.Warnf("Storage update returned status %d", resp.StatusCode)
	} else {
		h.logger.Info("User storage updated successfully")
	}
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
