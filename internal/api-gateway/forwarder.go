package api_gateway

import (
	"bytes"
	"io"
	"net/http"
)

// RequestForwarder определяет интерфейс для пересылки HTTP-запросов
// к другим сервисам. Реализации должны уметь пересылать запросы
// по указанному URL с сохранением метода, тела и заголовков.
type RequestForwarder interface {
	ForwardRequest(url string, r *http.Request) (*http.Response, error)
}

// DefaultRequestForwarder стандартная реализация интерфейса RequestForwarder
// для пересылки HTTP-запросов между сервисами.
type DefaultRequestForwarder struct{}

// ForwardRequest выполняет пересылку HTTP-запроса к указанному URL.
// Копирует метод, тело запроса и заголовки из исходного запроса.
// Возвращает ответ от целевого сервиса или ошибку в случае неудачи.
func (f *DefaultRequestForwarder) ForwardRequest(url string, r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	req, err := http.NewRequest(r.Method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header = r.Header
	return client.Do(req)
}
