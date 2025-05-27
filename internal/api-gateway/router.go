package api_gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// NewRouter создает новый маршрутизатор Gin с настройкой всех конечных точек API
// handler - обработчик запросов, содержащий логику для всех API endpoints
// Возвращает сконфигурированный экземпляр gin.Engine
func NewRouter(handler *APIHandler) *gin.Engine {
	r := gin.Default()

	r.Use(ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/upload", handler.UploadHandler)
	r.POST("/analyze", handler.AnalyzeHandler)
	r.POST("/compare", handler.CompareHandler)

	return r
}
