package filestoringservice

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// InitDB инициализирует подключение к базе данных и создает таблицы, если они не существуют.
func InitDB() (*sql.DB, error) {
	connStr := "postgresql://file_user:160206@postgres:5432/file_storage_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %v", err)
	}

	// Создаем таблицу если она не существует
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("ошибка создания таблиц: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("проверка подключения не удалась: %v", err)
	}

	return db, nil
}

// createTables создает таблицу 'files' в базе данных, если она не существует.
func createTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS files (
        id SERIAL PRIMARY KEY,
        filename TEXT NOT NULL,
        numberofparagraphs INTEGER NOT NULL,
        numberofwords INTEGER NOT NULL,
        numberofsymbols INTEGER NOT NULL,
        hashcode TEXT NOT NULL UNIQUE
    )`

	_, err := db.Exec(query)
	return err
}

// CloseDB закрывает соединение с базой данных.
func CloseDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия соединения: %v", err)
	}
	return nil
}

// CheckFileExists проверяет, существует ли файл с указанным хешом в базе данных.
func CheckFileExists(db *sql.DB, fileHash string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM files WHERE hashcode = $1)",
		fileHash,
	).Scan(&exists)
	return exists, err
}

// InsertMetadata вставляет метаданные файла в таблицу 'files'.
func InsertMetadata(db *sql.DB, filename, filehash string, paragraphs, words, symbols int) error {
	query := `INSERT INTO files (filename, numberofparagraphs, numberofwords, numberofsymbols, hashcode) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query, filename, paragraphs, words, symbols, filehash)
	return err
}
