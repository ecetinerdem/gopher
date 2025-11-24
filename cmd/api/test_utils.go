package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ecetinerdem/gopherSocial/internal/auth"
	"github.com/ecetinerdem/gopherSocial/internal/store"
	"github.com/ecetinerdem/gopherSocial/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()
	testAuthenticator := auth.NewMockAuthenticator()
	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCache,
		authenticator: testAuthenticator,
		config:        cfg,
	}
}

func executeRequest(r *http.Request, mux http.Handler) httptest.ResponseRecorder {

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, r)

	return *recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected the response code %d but received %d", expected, actual)
	}
}
