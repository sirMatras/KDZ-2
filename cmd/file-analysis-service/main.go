package main

import (
	filestoringservice "KDZ/internal/file-storing-service"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	"log"

	_ "KDZ/docs"
	fileanalysisservice "KDZ/internal/file-analysis-service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
)

// @title File Analysis API
// @version 1.0
// @description API для анализа файлов
// @contact.name API Support
// @contact.url https://www.example.com/support
// @contact.email support@example.com
// @host localhost:8082
// @BasePath /
func main() {
	db, err := filestoringservice.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer filestoringservice.CloseDB(db)

	router := gin.Default()

	router.POST("/analyze", fileanalysisservice.AnalyzeFile)
	router.POST("/compare", func(c *gin.Context) {
		fileanalysisservice.CompareFileByHash(c, db)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := ":8082"
	fmt.Printf("\nFile Analysis Service запущен на порту %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
