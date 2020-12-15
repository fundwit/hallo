package serveHttp

import (
	bytes2 "bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
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

func TestAccountHandler_createUser(it *testing.T) {
	it.Run("should create failed in invalid conditions", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

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
		accountHandler.RegisterRoutes(engine.Group("/accounts"))

		// --- bad body ---
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader("xxx"))
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse := w.Result()
		defer httpResponse.Body.Close()
		body, _ := ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
		responseBody, err := json.Marshal(gin.H{"error": "bad request body"})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(responseBody), string(body))

		registerToken := uuid.New().String()
		registerEmail := uuid.New().String() + "@test.fundwit.com"

		// --- register token is empty ---
		requestBody, err := json.Marshal(AccountCreateForm{
			Name: uuid.New().String(), Email: registerEmail, Secret: uuid.New().String(), RegisterToken: ""})
		req = httptest.NewRequest(http.MethodPost, "/accounts", bytes2.NewReader(requestBody))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
		responseBody, err = json.Marshal(gin.H{"error": "bad request body"})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(responseBody), string(body))

		// create register token
		auth.RegisterTokenCache.Set(registerEmail, registerToken, cache.DefaultExpiration)

		// --- register token is not match ---
		requestBody, err = json.Marshal(AccountCreateForm{
			Name: uuid.New().String(), Email: registerEmail, Secret: uuid.New().String(), RegisterToken: registerToken + ".bad"})
		req = httptest.NewRequest(http.MethodPost, "/accounts", bytes2.NewReader(requestBody))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
		responseBody, err = json.Marshal(gin.H{"error": (&domain.ErrRegisterTokenInvalid{}).Error()})
		if err != nil {
			panic(err)
		}
		assert.JSONEq(t, string(responseBody), string(body))

		// --- success ---
		createFrom := AccountCreateForm{Name: uuid.New().String(), Email: registerEmail, Secret: uuid.New().String(), RegisterToken: registerToken}
		requestBody, err = json.Marshal(createFrom)
		req = httptest.NewRequest(http.MethodPost, "/accounts", bytes2.NewReader(requestBody))
		w = httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		httpResponse = w.Result()
		defer httpResponse.Body.Close()
		body, _ = ioutil.ReadAll(httpResponse.Body)

		// assertion
		assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
		bodyJson := map[string]entity.Account{}
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, createFrom.Name, bodyJson["user"].Name)
		assert.Equal(t, createFrom.Email, bodyJson["user"].Email)
		assert.NotNil(t, bodyJson["user"].Id)
		assert.NotNil(t, bodyJson["user"].CreateTime)
		assert.NotNil(t, bodyJson["user"].LastUpdateTime)
		// assert token has been consumed
		_, found := auth.RegisterTokenCache.Get(registerEmail)
		assert.False(t, found)
	})
}
