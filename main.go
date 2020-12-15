package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hallo/bootstrap"
	"hallo/dataSource"
	"hallo/domain"
	"hallo/meta"
	"hallo/serveHttp"
	"hallo/service/auth"
	"hallo/util"
	"log"
)

func main() {
	ds, err := new(dataSource.DataSource).Start()
	if err != nil {
		panic(err)
	}
	defer ds.Stop()

	accountRepository := &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database}
	accountManager := &domain.AccountManagerImpl{
		AccountRepository:          accountRepository,
		IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
		InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
	}

	sessionHandler := serveHttp.SessionHandler{AccountManager: accountManager}
	accountHandler := serveHttp.AccountHandler{
		AccountManager:    accountManager,
		AccountRepository: accountRepository,
	}
	registryHandler := serveHttp.RegistryHandler{AccountRepository: accountRepository}

	_, err = bootstrap.CreateInitialAccount(accountManager, accountRepository)
	if err != nil {
		panic(fmt.Errorf("failed to check and prepare default admin account. %w", err))
	}

	engine := gin.Default()
	engine.Use(auth.AuthenticateByToken())

	meta.Routes(engine.Group("/"))
	sessionHandler.RegisterRoutes(engine.Group("/sessions"))
	accountHandler.RegisterRoutes(engine.Group("/accounts"))
	registryHandler.RegisterRoutes(engine.Group("/registry"))

	log.Println("service start")
	err = engine.Run(":80")
	if err != nil {
		panic(err)
	}
}
