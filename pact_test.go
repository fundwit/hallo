package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/gomock"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"github.com/patrickmn/go-cache"
	"hallo/domain"
	"hallo/domain/entity"
	"hallo/meta"
	"hallo/serveHttp"
	"hallo/service/auth"
	"os"
	"testing"
	"time"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)
var port, _ = utils.GetFreePort()

var mockAccountManager *domain.MockAccountManager
var mockAccountRepository *domain.MockAccountRepository

// The Provider verification
func TestPactProvider(t *testing.T) {
	// 启动 Provider 服务
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAccountManager = domain.NewMockAccountManager(mockCtl)
	mockAccountRepository = domain.NewMockAccountRepository(mockCtl)

	go startInstrumentedProvider()

	pact := dsl.Pact{
		Provider:                 meta.GetServiceInfo().ServiceName,
		LogDir:                   logDir,
		PactDir:                  pactDir,
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}

	// Verify the Provider - Tag-based Published Pacts for any known consumers
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://127.0.0.1:%d", port),
		//Tags: []string{"master"},  // consumer tags
		FailIfNoPactsFound: false,
		// Use this if you want to test without the Pact Broker
		// PactURLs:  []string{filepath.FromSlash(fmt.Sprintf("%s/goadminservice-gouserservice.json", os.Getenv("PACT_DIR")))},
		BrokerURL:                  "https://pact.fundwit.com",
		BrokerUsername:             os.Getenv("PACT_BROKER_USERNAME"),
		BrokerPassword:             os.Getenv("PACT_BROKER_PASSWORD"),
		PublishVerificationResults: true,
		ProviderVersion:            "1.0.0",
		StateHandlers:              stateHandlers,
		// RequestFilter: fixBearerToken,
	})

	if err != nil {
		t.Fatal(err)
	}
}

//// Simulates the need to set a time-bound authorization token,
//// such as an OAuth bearer toke
//func fixBearerToken(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
//		// Only set the correct bearer token, if one was provided in the first place
//		if request.Header.Get("Authorization") != "" {
//			request.Header.Set("Authorization", "xxx")
//		}
//		next.ServeHTTP(writer, request)
//	})
//}

var stateHandlers = types.StateHandlers{
	"email [available@test.fundwit.com] not occupied": func() error {
		mockAccountRepository.EXPECT().IsEmailOccupied("available@test.fundwit.com").Return(false, nil)
		return nil
	},
	"email [occupied@test.fundwit.com] occupied": func() error {
		// create a user to occupy the email
		mockAccountRepository.EXPECT().IsEmailOccupied("occupied@test.fundwit.com").Return(true, nil)
		return nil
	},
	"account name [Ann] not occupied": func() error {
		mockAccountRepository.EXPECT().IsAccountNameOccupied("Ann").Return(false, nil)
		return nil
	},
	"account name [Bob] occupied": func() error {
		mockAccountRepository.EXPECT().IsAccountNameOccupied("Bob").Return(true, nil)
		return nil
	},

	"success login with credential [Ann, correctSecret]": func() error {
		mockAccountManager.EXPECT().AuthenticateInternalIdentity("Ann", "correctSecret").Return(
			&entity.Account{Name: "Ann", Email: "ann@test.fundwit.com", Id: 123, CreateTime: time.Now(), LastUpdateTime: time.Now()}, nil)
		return nil
	},
	"failed login with credential [Ann, badSecret]": func() error {
		mockAccountManager.EXPECT().AuthenticateInternalIdentity("Ann", "badSecret").Return(nil, &domain.AccountAuthenticationFailure{})
		return nil
	},
	"success logout": func() error {
		securityContext := &auth.SecurityContext{Token: "correctToken", Principal: auth.Principal{Name: "Ann"}}
		auth.TokenCache.Set("correctToken", securityContext, cache.DefaultExpiration)
		return nil
	},
	"failed logout": func() error {
		auth.TokenCache.Delete("badToken")
		return nil
	},
	"already login for session info": func() error {
		securityContext := &auth.SecurityContext{Token: "correctToken", Principal: auth.Principal{Name: "Ann"}}
		auth.TokenCache.Set("correctToken", securityContext, cache.DefaultExpiration)
		return nil
	},
	"un-login for session info": func() error {
		// clean cache
		auth.TokenCache = cache.New(24*time.Hour, 1*time.Minute)
		return nil
	},
	"session info - bad token": func() error {
		// clean cache
		auth.TokenCache = cache.New(24*time.Hour, 1*time.Minute)
		return nil
	},

	"sign up success with parameters [Ann, email-sign-up@test.fundwit.com, correct_register_token, correctSecret]": func() error {
		auth.RegisterTokenCache.Set("email-sign-up@test.fundwit.com", "correct_register_token", cache.DefaultExpiration)
		mockAccountManager.EXPECT().CreateAccount(entity.EmailAccountCreateRequest{Name: "Ann", Secret: "correctSecret", Email: "email-sign-up@test.fundwit.com"}).
			Return(&entity.Account{Name: "Ann", Email: "email-sign-up@test.fundwit.com", Id: 123, CreateTime: time.Now(), LastUpdateTime: time.Now()}, nil)
		return nil
	},
	"failed sign up": func() error {
		// create a user
		return nil
	},
}

// Starts the provider API with hooks for provider states.
// This essentially mirrors the main.go file, with extra routes added.
func startInstrumentedProvider() {
	sessionHandler := serveHttp.SessionHandler{AccountManager: mockAccountManager}
	accountHandler := serveHttp.AccountHandler{
		AccountManager:    mockAccountManager,
		AccountRepository: mockAccountRepository,
	}
	registryHandler := serveHttp.RegistryHandler{AccountRepository: mockAccountRepository}

	engine := gin.Default()
	engine.Use(auth.AuthenticateByToken())

	meta.Routes(engine.Group("/"))
	sessionHandler.RegisterRoutes(engine.Group("/sessions"))
	accountHandler.RegisterRoutes(engine.Group("/accounts"))
	registryHandler.RegisterRoutes(engine.Group("/registry"))

	engine.Run(fmt.Sprintf(":%d", port))
}

//------------------
// mocks
//------------------
