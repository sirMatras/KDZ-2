package file_analysis_service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

// FileStats содержит статистику анализа файла:
// Paragraphs - количество параграфов
// Words - количество слов
// Symbols - количество символов
type FileStats struct {
	Paragraphs int `json:"paragraphs"`
	Words      int `json:"words"`
	Symbols    int `json:"symbols"`
}

// AnalyzeFile анализирует файл и возвращает статистику
// @Summary Анализ файла
// @Description Анализирует файл и возвращает статистику
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} FileStats "Успешный ответ"
// @Failure 400 {object} map[string]interface{} "Ошибка запроса"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /analyze [post]
func AnalyzeFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный запрос"})
		return
	}

	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения файла"})
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка открытия файла"})
		return
	}
	defer f.Close()

	content := ""
	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if err != nil && err.Error() != "EOF" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения файла"})
			return
		}
		if n == 0 {
			break
		}
		content += string(buf[:n])
	}

	text := content
	paragraphs := len(strings.Split(text, "\n\n"))
	words := len(strings.Fields(text))
	symbols := len(text)

	stats := FileStats{
		Paragraphs: paragraphs,
		Words:      words,
		Symbols:    symbols,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "успешно",
		"stats":   stats,
	})
}
