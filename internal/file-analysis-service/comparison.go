package file_analysis_service

import (
	filestoringservice "KDZ/internal/file-storing-service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CompareFileByHash - сравнивает файл с файлами в базе данных по хешу
// @Summary Сравнение файлов по содержимому
// @Description Принимает файл, вычисляет его хеш и проверяет, существует ли файл с таким же хешем в базе данных
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для сравнения"
// @Success 200 {object} map[string]interface{} "Файл найден или нет"
// @Failure 400 {object} map[string]interface{} "Ошибка запроса"
// @Failure 404 {object} map[string]interface{} "Файл не найден"
// @Router /compare [post]
func CompareFileByHash(c *gin.Context, db *sql.DB) {
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

	hash, err := filestoringservice.CalculateFileHash(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка вычисления хеша файла"})
		return
	}

	exists, err := filestoringservice.CheckFileExists(db, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
		return
	}

	if exists {
		c.JSON(http.StatusOK, gin.H{
			"message": "файл существует на 100%",
			"hash":    hash,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "файл не найден. Нет совпадений.",
			"hash":    hash,
		})
	}
}
