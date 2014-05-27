package perfect

import (
	"net/http"
)

type mockResponse struct {
	Headers http.Header
	Status  int
	Data    []byte
}

func (m *mockResponse) Header() http.Header {
	return m.Headers
}

func NewMockResponse() *mockResponse {
	return &mockResponse{
		Headers: http.Header{},
	}
}

func (m *mockResponse) Write(data []byte) (bytes_written int, err error) {
	m.Data = append(m.Data, data...)
	err = nil
	bytes_written = len(data)

	return
}

func (m *mockResponse) WriteHeader(status int) {
	m.Status = status
}
