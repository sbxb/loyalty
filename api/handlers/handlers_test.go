package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/loyalty/api/handlers"
	mw "github.com/sbxb/loyalty/api/middleware"
	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg config.Config

var _ = func() bool {
	// stackoverflow.com-recommended hack to parse testing flags before
	// application ones - prevents test failure with an error:
	// "flag provided but not defined: -test.testlogfile"
	testing.Init()

	var err error
	if cfg, err = config.New(); err != nil {
		log.Fatal(err)
	}
	return true
}()

func TestUserRegister_NotValidInput(t *testing.T) {
	wantCode := 400
	tests := []string{
		"",
		"   ",
		"abc",
		"{}",
		"[]",
		`{"key": "value"}`,
		`{"login": "<login>", "password": "<password>", "Hash": "<data>"}`,
		`{"login": "<login>"}`,
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error
	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/user/register", urlHandler.UserRegister)

	for _, tt := range tests {
		t.Run("Register", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://"+cfg.ServerAddress+"/api/user/register", strings.NewReader(tt))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			//resBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, resp.StatusCode, wantCode)
			//t.Log(string(resBody))
		})
	}
}

func TestUserRegister_ValidInput(t *testing.T) {
	user := &models.User{
		Login:    "user",
		Password: "abcdef",
	}
	tests := []struct {
		wantCode      int
		user          *models.User
		hasUserCookie bool
	}{
		{
			wantCode:      200,
			user:          user,
			hasUserCookie: true,
		},
		{
			wantCode:      409,
			user:          user,
			hasUserCookie: false,
		},
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/user/register", urlHandler.UserRegister)

	for _, tt := range tests {
		t.Run("Register", func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.user)
			req := httptest.NewRequest(
				http.MethodPost,
				"http://"+cfg.ServerAddress+"/api/user/register",
				bytes.NewReader(requestBody),
			)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			//resBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, resp.StatusCode, tt.wantCode)
			//t.Log(string(resBody))
			assert.Equal(t, checkCookie(resp, "user"), tt.hasUserCookie)
		})
	}
}

func TestUserLogin_NotValidInput(t *testing.T) {
	wantCode := 400
	tests := []string{
		"",
		"   ",
		"abc",
		"{}",
		"[]",
		`{"key": "value"}`,
		`{"login": "<login>", "password": "<password>", "Hash": "<data>"}`,
		`{"login": "<login>"}`,
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error
	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/user/login", urlHandler.UserLogin)

	for _, tt := range tests {
		t.Run("Login", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://"+cfg.ServerAddress+"/api/user/login", strings.NewReader(tt))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			//resBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, resp.StatusCode, wantCode)
			//t.Log(string(resBody))
		})
	}
}

func TestUserLogin_ValidInput(t *testing.T) {
	tests := []struct {
		wantCode      int
		user          *models.User
		hasUserCookie bool
	}{
		{
			wantCode: 200,
			user: &models.User{
				Login:    "user",
				Password: "abcdef",
			},
			hasUserCookie: true,
		},
		{
			wantCode: 401,
			user: &models.User{
				Login:    "nonexistentuser",
				Password: "abcdef",
			},
			hasUserCookie: false,
		},
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error
	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/user/login", urlHandler.UserLogin)

	// add the first user
	user := &models.User{
		Login: "user",
		Hash:  "$2a$10$2V0TfI3A/Win8OI5Q.U1gOjffxfBxX9bLUa7Zheo3jKOaxAzwEDYa",
	}
	err := store.AddUser(context.Background(), user)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run("Login", func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.user)
			req := httptest.NewRequest(
				http.MethodPost,
				"http://"+cfg.ServerAddress+"/api/user/login",
				bytes.NewReader(requestBody),
			)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			//resBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, resp.StatusCode, tt.wantCode)
			//t.Log(string(resBody))
			assert.Equal(t, checkCookie(resp, "user"), tt.hasUserCookie)
		})
	}
}

func TestUserPostOrder_ValidInput(t *testing.T) {
	tests := []struct {
		wantCode    int
		orderNumber string
	}{
		{
			wantCode:    202,
			orderNumber: "12345678903",
		},
		{
			wantCode:    200,
			orderNumber: "12345678903",
		},
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error
	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.With(mw.AuthMW).Post("/api/user/orders", urlHandler.UserPostOrder)

	// add the first user
	user := &models.User{
		Login: "user",
		Hash:  "$2a$10$2V0TfI3A/Win8OI5Q.U1gOjffxfBxX9bLUa7Zheo3jKOaxAzwEDYa",
	}
	err := store.AddUser(context.Background(), user)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run("Add Order", func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"http://"+cfg.ServerAddress+"/api/user/orders",
				strings.NewReader(tt.orderNumber),
			)
			cookie := http.Cookie{
				Name:    "user",
				Value:   "47dd2e0ab2fc35c8bd9d56847a904a3ca3cf166c5c1ffb1b4f8a87337b134ebcNoWIWEqALp/7+sHKF7Dq5/mpSvupLXgkpj6TDcdXeu/almqW9dykJ6hVZQGx/nPnhlFv",
				Expires: time.Now().Add(1 * time.Hour),
			}
			req.AddCookie(&cookie)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, tt.wantCode)
			t.Log(resp.StatusCode)
			store.DumpOrder()
			store.DumpBalance()
		})
	}

}

func checkCookie(resp *http.Response, key string) bool {
	for _, c := range resp.Cookies() {
		if c.Name == key {
			fmt.Println(c.Value)
			return true
		}
	}
	return false
}
