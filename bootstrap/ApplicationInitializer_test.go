package bootstrap

import (
	"github.com/stretchr/testify/assert"
	"hallo/domain"
	"hallo/testinfra"
	"hallo/util"
	"testing"
)

func TestCreateInitialAccount(it *testing.T) {
	it.Run("should create default admin account correctly", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		accountManager := domain.AccountManagerImpl{
			AccountRepository:          &domain.DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
			IdentityBindingRepository:  &domain.DatabaseIdentityBindingRepository{Database: ds.Database},
			InternalIdentityRepository: &domain.DatabaseInternalIdentityRepository{Database: ds.Database},
		}

		account, err := CreateInitialAccount(&accountManager, accountManager.AccountRepository)
		assert.Nil(t, err)
		assert.Equal(t, "admin", account.Name)

		account, err = accountManager.AuthenticateInternalIdentity("admin", "admin123")
		assert.Nil(t, err)
		assert.Equal(t, "admin", account.Name)

		// case 2
		account, err = CreateInitialAccount(&accountManager, accountManager.AccountRepository)
		assert.Nil(t, err)
		assert.Nil(t, account)
	})
}
