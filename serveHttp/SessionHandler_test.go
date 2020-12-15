package serveHttp

import (
	bytes2 "bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"hallo/domain"
	"hallo/domain/entity"
	"hallo/service/auth"
	"hallo/testinfra"
	"hallo/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSessionHandler_newSession(it *testing.T) {
	it.Run("should login failed when binding failed", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}

		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		req := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader("xxx"))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()
		body, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
		token := httpResponse.Header.Get("Authentication")
		assert.Equal(t, "", token)
		responseBody, err := json.Marshal(gin.H{"error": "bad request body"})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(responseBody), string(body))
	})

	it.Run("should login failed when account not exist or secret is not match", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}
		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		accountName := uuid.New().String()
		accountSecret := uuid.New().String()

		requestBody, err := json.Marshal(LoginRequest{Name: accountName, Secret: accountSecret})
		if err != nil {
			panic(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/sessions", bytes2.NewReader(requestBody))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()
		body, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)
		token := httpResponse.Header.Get("Authentication")
		assert.Equal(t, "", token)
		responseBody, err := json.Marshal(gin.H{"error": "account not exist or secret is not match"})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(responseBody), string(body))

		// ----- bad secret -----
		// create account
		_, err = sessionHandler.AccountManager.CreateAccount(
			entity.EmailAccountCreateRequest{Name: accountName, Email: accountName + "@test.fundwit.com", Secret: accountSecret})

		requestBody, err = json.Marshal(LoginRequest{Name: accountName, Secret: accountSecret + "bad"})
		if err != nil {
			panic(err)
		}
		req = httptest.NewRequest(http.MethodPost, "/sessions", bytes2.NewBuffer(requestBody))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)
		token = httpResponse.Header.Get("Authentication")
		assert.Equal(t, "", token)
		responseBody, err = json.Marshal(gin.H{"error": "account not exist or secret is not match"})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(responseBody), string(body))
	})

	it.Run("should login success", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}
		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		// create account
		accountName := uuid.New().String()
		accountSecret := uuid.New().String()
		_, err := sessionHandler.AccountManager.CreateAccount(
			entity.EmailAccountCreateRequest{Name: accountName, Email: accountName + "@test.fundwit.com", Secret: accountSecret})

		requestBody, err := json.Marshal(LoginRequest{Name: accountName, Secret: accountSecret})
		if err != nil {
			panic(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/sessions", bytes2.NewReader(requestBody))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()
		responseBody, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		token := httpResponse.Header.Get("Authentication")
		assert.NotNil(t, token)
		wantedBody, err := json.Marshal(gin.H{"token": token, "principal": gin.H{"name": accountName}})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(wantedBody), string(responseBody))

		// token was cached
		sc, found := auth.TokenCache.Get(token)
		assert.True(t, found)
		want := &auth.SecurityContext{Token: token, Principal: auth.Principal{Name: accountName}}
		assert.Equal(t, want, sc)
	})
}

func TestSessionHandler_deleteSession(it *testing.T) {
	it.Run("should delete success failed when session not exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}

		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		// --- logout without token ---
		req := httptest.NewRequest(http.MethodDelete, "/sessions", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		// assertion
		assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
	})

	it.Run("should delete success when session is exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}

		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		// create account
		accountName := uuid.New().String()
		accountSecret := uuid.New().String()
		_, err := sessionHandler.AccountManager.CreateAccount(
			entity.EmailAccountCreateRequest{Name: accountName, Email: accountName + "@test.fundwit.com", Secret: accountSecret})

		// login account
		requestBody, err := json.Marshal(LoginRequest{Name: accountName, Secret: accountSecret})
		if err != nil {
			panic(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/sessions", bytes2.NewReader(requestBody))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		// get token
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		token := httpResponse.Header.Get("Authentication")
		assert.NotNil(t, token)
		sc, found := auth.TokenCache.Get(token)
		assert.True(t, found)
		want := &auth.SecurityContext{Token: token, Principal: auth.Principal{Name: accountName}}
		assert.Equal(t, want, sc)

		// --- logout with bad token ---
		req = httptest.NewRequest(http.MethodDelete, "/sessions", nil)
		req.Header.Set("Authorization", "bearer "+token+".bad")
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse = w.Result()
		defer httpResponse.Body.Close()

		// assertion
		assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
		// token still exist cleaned
		sc, found = auth.TokenCache.Get(token)
		assert.True(t, found)

		// --- logout with correct token ---
		req = httptest.NewRequest(http.MethodDelete, "/sessions", nil)
		req.Header.Set("Authorization", "bearer "+token)
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse = w.Result()
		defer httpResponse.Body.Close()

		// assertion
		assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
		// token has been cleaned
		sc, found = auth.TokenCache.Get(token)
		assert.False(t, found)
	})
}

func TestSessionHandler_currentSession(it *testing.T) {
	it.Run("should failed when session not exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}

		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		// --- logout without token ---
		req := httptest.NewRequest(http.MethodGet, "/sessions/me", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		// assertion
		assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)
	})

	it.Run("should get correct result when session exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		sessionHandler := SessionHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
		}

		engine := gin.Default()
		sessionHandler.RegisterRoutes(engine.Group("/sessions"))

		// create account
		accountName := uuid.New().String()
		accountSecret := uuid.New().String()
		_, err := sessionHandler.AccountManager.CreateAccount(
			entity.EmailAccountCreateRequest{Name: accountName, Email: accountName + "@test.fundwit.com", Secret: accountSecret})

		// login account
		requestBody, err := json.Marshal(LoginRequest{Name: accountName, Secret: accountSecret})
		if err != nil {
			panic(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/sessions", bytes2.NewReader(requestBody))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		// get token
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		token := httpResponse.Header.Get("Authentication")
		assert.NotNil(t, token)
		sc, found := auth.TokenCache.Get(token)
		assert.True(t, found)
		want := &auth.SecurityContext{Token: token, Principal: auth.Principal{Name: accountName}}
		assert.Equal(t, want, sc)

		// --- get session with bad token ---
		req = httptest.NewRequest(http.MethodGet, "/sessions/me", nil)
		req.Header.Set("Authorization", "bearer "+token+".bad")
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse = w.Result()
		defer httpResponse.Body.Close()

		// assertion
		assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)

		// --- get session with correct token ---
		req = httptest.NewRequest(http.MethodGet, "/sessions/me", nil)
		req.Header.Set("Authorization", "bearer "+token)
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		responseBody, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		wantedBody, err := json.Marshal(gin.H{"token": token, "principal": gin.H{"name": accountName}})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(wantedBody), string(responseBody))
	})
}
