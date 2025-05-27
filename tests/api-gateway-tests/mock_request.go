package api_gateway_tests

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockRequestForwarder struct {
	mock.Mock
}

func (m *MockRequestForwarder) ForwardRequest(url string, req *http.Request) (*http.Response, error) {
	args := m.Called(url, req)
	return args.Get(0).(*http.Response), args.Error(1)
}
