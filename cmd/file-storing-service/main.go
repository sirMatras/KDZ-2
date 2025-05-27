package main

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	"log"
	"os"

	_ "KDZ/docs"
	filestoringservice "KDZ/internal/file-storing-service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
)

// @title File Storage API
// @version 1.0
// @description API для хранения и загрузки файлов
// @contact.name API Support
// @contact.url https://www.example.com/support
// @contact.email support@example.com
// @host localhost:8081
// @BasePath /
func main() {
	db, err := filestoringservice.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer filestoringservice.CloseDB(db)

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal("Ошибка создания папки:", err)
	}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Используем правильный импорт

	router.POST("/upload", func(c *gin.Context) {
		filestoringservice.SaveFile(c, db)
	})

	port := ":8081"
	fmt.Printf("\nFile Storing Service запущен на порту %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
