package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KazikovAP/merch_store/internal/api/http/middleware"
	"github.com/stretchr/testify/assert"
)

// Тест на запись ошибки.
func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	middleware.WriteError(w, http.StatusBadRequest, "test error message")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test error message", response["error"])
}
