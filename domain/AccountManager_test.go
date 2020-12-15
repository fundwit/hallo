package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"hallo/domain/entity"
	"hallo/testinfra"
	"hallo/util"
	"testing"
)

func TestAccountManager_CreateAccount(it *testing.T) {
	it.Run("should create account successful when no conflict", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		accountManager := AccountManagerImpl{
			AccountRepository:          &DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
			IdentityBindingRepository:  &DatabaseIdentityBindingRepository{Database: ds.Database},
			InternalIdentityRepository: &DatabaseInternalIdentityRepository{Database: ds.Database},
		}

		accountName := uuid.New().String()
		accountSecret := uuid.New().String()

		found, err := accountManager.AccountRepository.IsAccountNameOccupied(accountName)
		assert.False(t, found, err)
		assert.Equal(t, nil, err)

		// create user
		account, err := accountManager.CreateAccount(entity.EmailAccountCreateRequest{
			Name: accountName, Secret: accountSecret, Email: accountName + "@test.fundwit.com",
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, accountName, account.Name)

		// verify: account is created
		found, err = accountManager.AccountRepository.IsAccountNameOccupied(accountName)
		assert.True(t, found, err)
		assert.Equal(t, nil, err)

		// verify internalIdentity record is created: verified when test authenticate

		// TODO verify: identityBinding record is created
	})

	it.Run("should create account failed when conflict", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		accountManager := AccountManagerImpl{
			AccountRepository:          &DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
			IdentityBindingRepository:  &DatabaseIdentityBindingRepository{Database: ds.Database},
			InternalIdentityRepository: &DatabaseInternalIdentityRepository{Database: ds.Database},
		}

		accountName := uuid.New().String()
		accountSecret := uuid.New().String()

		// create user
		account, err := accountManager.CreateAccount(entity.EmailAccountCreateRequest{
			Name: accountName, Secret: accountSecret, Email: accountName + "@test.fundwit.com",
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, accountName, account.Name)

		account, err = accountManager.CreateAccount(entity.EmailAccountCreateRequest{
			Name: accountName, Secret: accountSecret, Email: accountName + "1@test.fundwit.com",
		})
		assert.Equal(t, &AccountNameIsOccupied{}, err)
		assert.Nil(t, account)

		account, err = accountManager.CreateAccount(entity.EmailAccountCreateRequest{
			Name: accountName + "1", Secret: accountSecret, Email: accountName + "@test.fundwit.com",
		})
		assert.Equal(t, &AccountEmailIsOccupied{}, err)
		assert.Nil(t, account)
	})

}

func TestAccountManager_AuthenticateInternalIdentity(it *testing.T) {
	it.Run("should authenticate failed when account is not exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		accountManager := AccountManagerImpl{
			AccountRepository:          &DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
			IdentityBindingRepository:  &DatabaseIdentityBindingRepository{Database: ds.Database},
			InternalIdentityRepository: &DatabaseInternalIdentityRepository{Database: ds.Database},
		}
		accountName := uuid.New().String()
		accountSecret := uuid.New().String()

		account, err := accountManager.AuthenticateInternalIdentity(accountName, accountSecret)
		assert.Nil(t, account)
		assert.True(t, gorm.IsRecordNotFoundError(err))
	})

	it.Run("should authenticate correct when account is exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		accountManager := AccountManagerImpl{
			AccountRepository:          &DatabaseAccountRepository{IdWorker: util.DefaultIdWorker, Database: ds.Database},
			IdentityBindingRepository:  &DatabaseIdentityBindingRepository{Database: ds.Database},
			InternalIdentityRepository: &DatabaseInternalIdentityRepository{Database: ds.Database},
		}

		// create user
		accountName := uuid.New().String()
		accountSecret := uuid.New().String()
		_, err := accountManager.CreateAccount(entity.EmailAccountCreateRequest{
			Name: accountName, Secret: accountSecret, Email: accountName + "@test.fundwit.com",
		})
		if err != nil {
			panic(err)
		}

		account, err := accountManager.AuthenticateInternalIdentity(accountName, accountSecret)
		assert.Equal(t, accountName, account.Name)
		assert.Nil(t, err)

		account, err = accountManager.AuthenticateInternalIdentity(accountName, accountSecret+"bad")
		assert.Nil(t, account)
		assert.Equal(t, &AccountAuthenticationFailure{}, err)
	})
}
