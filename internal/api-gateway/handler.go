package api_gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
)

// APIHandler обрабатывает HTTP-запросы и перенаправляет их в соответствующие сервисы
type APIHandler struct {
	Forwarder RequestForwarder
}

// handleRequest общая функция для обработки и перенаправления запросов
// url - адрес целевого сервиса
// c - контекст Gin для обработки HTTP-запроса/ответа
func (h *APIHandler) handleRequest(url string, c *gin.Context) {
	resp, err := h.Forwarder.ForwardRequest(url, c.Request)
	if err != nil {
		log.Println("Error forwarding request:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding JSON:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(resp.StatusCode, result)
}

// UploadHandler обрабатывает запросы на загрузку файлов,
// перенаправляя их в file-storing-service
func (h *APIHandler) UploadHandler(c *gin.Context) {
	h.handleRequest("http://file-storing-service:8081/upload", c)
}

// AnalyzeHandler обрабатывает запросы на анализ файлов,
// перенаправляя их в file-analysis-service
func (h *APIHandler) AnalyzeHandler(c *gin.Context) {
	h.handleRequest("http://file-analysis-service:8082/analyze", c)
}

// CompareHandler обрабатывает запросы на сравнение файлов,
// поддерживает как JSON-ответы, так и произвольные данные
func (h *APIHandler) CompareHandler(c *gin.Context) {
	resp, err := h.Forwarder.ForwardRequest("http://file-analysis-service:8082/compare", c.Request)
	if err != nil {
		log.Println("Error forwarding compare request:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err == nil {
		c.JSON(resp.StatusCode, result)
	} else {
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = http.DetectContentType(bodyBytes)
		}
		c.Data(resp.StatusCode, contentType, bodyBytes)
	}
}
