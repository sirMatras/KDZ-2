package filestoringservice

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Загрузка файла
// @Description Загружает файл и сохраняет его метаданные
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]interface{} "Успешный ответ"
// @Failure 400 {object} map[string]interface{} "Ошибка запроса"
// @Router /upload [post]
func SaveFile(c *gin.Context, db *sql.DB) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный запрос"})
		return
	}

	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения"})
		return
	}

	hash, err := CalculateFileHash(filePath) // Передаем строку с путем
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка хеширования"})
		return
	}

	exists, err := CheckFileExists(db, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "файл существует"})
		return
	}

	content, _ := os.ReadFile(filePath)
	text := string(content)

	paragraphs := len(strings.Split(text, "\n\n"))
	words := len(strings.Fields(text))
	symbols := len(text)

	if err := InsertMetadata(db, file.Filename, hash, paragraphs, words, symbols); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения метаданных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "успешно",
		"hash":    hash,
		"stats": gin.H{
			"paragraphs": paragraphs,
			"words":      words,
			"symbols":    symbols,
		},
	})
}

// CalculateFileHash вычисляет SHA256 хеш для файла по пути filePath
func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("ошибка вычисления хеша: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
