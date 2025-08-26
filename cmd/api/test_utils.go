package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/szoumoc/social/internal/auth"
	"github.com/szoumoc/social/internal/store"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d actual %d", expected, actual)
	}
}
