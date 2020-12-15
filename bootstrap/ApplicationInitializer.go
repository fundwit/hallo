package bootstrap

import (
	"hallo/domain"
	"hallo/domain/entity"
	"log"
	"os"
)

func CreateInitialAccount(accountManager domain.AccountManager, repository domain.AccountRepository) (createdAccount *entity.Account, err error) {
	accountName := "admin"
	accountSecret := os.Getenv("ADMIN_SECRET")
	if accountSecret == "" {
		accountSecret = "admin123"
	}

	count, err := repository.Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		log.Println("[INIT.ACCOUNT] some accounts are existed, default admin account will not be create")
		return nil, nil
	}

	account, err := accountManager.CreateAccount(entity.EmailAccountCreateRequest{
		Name:   accountName,
		Secret: accountSecret,
		Email:  "temp@test.fundwit.com",
	})

	if err == nil {
		log.Println("[INIT.ACCOUNT] default admin account has been created!")
	}

	return account, err
}
