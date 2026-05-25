package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsPreflightAcceptsArbitraryHeaders(t *testing.T) {
	router := newTestRouter()

	cases := []struct {
		name           string
		method         string
		path           string
		requestMethod  string
		requestHeaders string
	}{
		{"login content-type", http.MethodOptions, "/login", "POST", "Content-Type"},
		{"jogos custom header", http.MethodOptions, "/jogos", "POST", "X-Requested-With,Content-Type"},
		{"update authorization", http.MethodOptions, "/jogos/1", "PUT", "Authorization,Content-Type"},
		{"delete request", http.MethodOptions, "/jogos/1", "DELETE", "Authorization"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Origin", "http://localhost:3000")
			req.Header.Set("Access-Control-Request-Method", tc.requestMethod)
			req.Header.Set("Access-Control-Request-Headers", tc.requestHeaders)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
			assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
			assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
			assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
		})
	}
}

func TestCorsActualRequestSendsAllowOrigin(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/jogos", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
}
