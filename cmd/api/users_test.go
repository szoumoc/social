package main

import (
	"net/http"
	"testing"
)

// go test
func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		//check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)

	})
	t.Run("should allow authenticated requests", func(t *testing.T) {
		//check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
