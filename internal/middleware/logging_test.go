package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLoggerWritesStructuredRequestLog(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/auth/register", nil)
	req.Header.Set(requestIDHeader, "request-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "request-1", rec.Header().Get(requestIDHeader))
	require.Equal(t, 1, logs.Len())

	entry := logs.All()[0]
	assert.Equal(t, "http request", entry.Message)
	assert.Equal(t, http.MethodPost, entry.ContextMap()["method"])
	assert.Equal(t, "/auth/register", entry.ContextMap()["path"])
	assert.Equal(t, int64(http.StatusCreated), entry.ContextMap()["status"])
	assert.Equal(t, "request-1", entry.ContextMap()["request_id"])
	assert.NotNil(t, entry.ContextMap()["duration"])
}
