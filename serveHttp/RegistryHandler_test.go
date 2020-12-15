package serveHttp

import (
	bytes2 "bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"hallo/domain"
	"hallo/service/auth"
	"hallo/testinfra"
	"hallo/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegistryHandler_emailOccupied(it *testing.T) {
	it.Run("should response correctly for email occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		registryHandler := RegistryHandler{
			AccountRepository: &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
		}
		accountRepository := &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database}
		accountHandler := AccountHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          accountRepository,
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
			AccountRepository: accountRepository,
		}

		engine := gin.Default()
		registryHandler.RegisterRoutes(engine.Group("/registry"))
		accountHandler.RegisterRoutes(engine.Group("/accounts"))

		registerEmail := uuid.New().String() + "@test.fundwit.com"
		// --- email not occupied ---
		req := httptest.NewRequest(http.MethodPost, "/registry/emails", strings.NewReader("{\"email\": \""+registerEmail+"\"}"))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()
		body, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		wantedBody, err := json.Marshal(gin.H{"email": registerEmail, "occupied": false})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(wantedBody), string(body))

		// --- email already occupied ---
		// create account and occupied the email
		registerToken := uuid.New().String()
		auth.RegisterTokenCache.Set(registerEmail, registerToken, cache.DefaultExpiration)
		requestBody, err := json.Marshal(AccountCreateForm{
			Name: uuid.New().String(), Email: registerEmail, Secret: uuid.New().String(), RegisterToken: registerToken})
		req = httptest.NewRequest(http.MethodPost, "/accounts", bytes2.NewReader(requestBody))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)
		// assertion
		assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)

		req = httptest.NewRequest(http.MethodPost, "/registry/emails", strings.NewReader("{\"email\": \""+registerEmail+"\"}"))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		wantedBody, err = json.Marshal(gin.H{"email": registerEmail, "occupied": true})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(wantedBody), string(body))
	})

	it.Run("should response correctly for account name occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		registryHandler := RegistryHandler{
			AccountRepository: &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
		}
		accountRepository := &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database}
		accountHandler := AccountHandler{
			AccountManager: &domain.AccountManagerImpl{
				AccountRepository:          accountRepository,
				IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
				InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
			},
			AccountRepository: accountRepository,
		}

		engine := gin.Default()
		registryHandler.RegisterRoutes(engine.Group("/registry"))
		accountHandler.RegisterRoutes(engine.Group("/accounts"))

		registerName := uuid.New().String()
		// --- email not occupied ---
		req := httptest.NewRequest(http.MethodPost, "/registry/names", strings.NewReader("{\"name\": \""+registerName+"\"}"))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()
		body, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		wantedBody, err := json.Marshal(gin.H{"name": registerName, "occupied": false})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(wantedBody), string(body))

		// --- name already occupied ---
		// create account and occupied the email
		registerToken := uuid.New().String()
		registerEmail := uuid.New().String() + "@test.fundwit.com"
		auth.RegisterTokenCache.Set(registerEmail, registerToken, cache.DefaultExpiration)
		requestBody, err := json.Marshal(AccountCreateForm{
			Name: registerName, Email: registerEmail, Secret: uuid.New().String(), RegisterToken: registerToken})
		req = httptest.NewRequest(http.MethodPost, "/accounts", bytes2.NewReader(requestBody))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)
		// assertion
		assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)

		req = httptest.NewRequest(http.MethodPost, "/registry/names", strings.NewReader("{\"name\": \""+registerName+"\"}"))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)
		assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		wantedBody, err = json.Marshal(gin.H{"name": registerName, "occupied": true})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(wantedBody), string(body))
	})
}
