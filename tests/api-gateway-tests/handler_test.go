package api_gateway_tests

import (
	"KDZ/internal/api-gateway"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadHandler(t *testing.T) {
	mockForwarder := new(MockRequestForwarder)
	apiHandler := &api_gateway.APIHandler{Forwarder: mockForwarder}

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "success"}`))),
	}
	mockForwarder.On("ForwardRequest", "http://file-storing-service:8081/upload", mock.Anything).Return(mockResponse, nil)

	router := gin.Default()
	router.POST("/upload", apiHandler.UploadHandler)

	req, _ := http.NewRequest(http.MethodPost, "/upload", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "success"}`, w.Body.String())

	mockForwarder.AssertExpectations(t)
}

func TestAnalyzeHandler(t *testing.T) {
	mockForwarder := new(MockRequestForwarder)
	apiHandler := &api_gateway.APIHandler{Forwarder: mockForwarder}

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"analysis": "done"}`))),
	}
	mockForwarder.On("ForwardRequest", "http://file-analysis-service:8082/analyze", mock.Anything).Return(mockResponse, nil)

	router := gin.Default()
	router.POST("/analyze", apiHandler.AnalyzeHandler)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"analysis": "done"}`, w.Body.String())

	mockForwarder.AssertExpectations(t)
}
