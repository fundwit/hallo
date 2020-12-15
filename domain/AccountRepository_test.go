package domain

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"hallo/domain/entity"
	"hallo/testinfra"
	"hallo/util"
	"testing"
	"time"
)

func TestDatabaseAccountRepository_Count(it *testing.T) {
	it.Run("should return correct count amount", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		count, err := repository.Count()
		assert.Equal(t, nil, err)
		assert.Equal(t, uint64(0), count)

		ds.Database.Save(entity.Account{
			Id:             111,
			Name:           "test-count-1",
			Email:          "test-count-1@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 111})
		count, err = repository.Count()
		assert.Equal(t, nil, err)
		assert.Equal(t, uint64(1), count)

		ds.Database.Save(entity.Account{
			Id:             222,
			Name:           "test-count-2",
			Email:          "test-count-2@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 222})
		count, err = repository.Count()
		assert.Equal(t, nil, err)
		assert.Equal(t, uint64(2), count)
	})

	it.Run("should return true if email is occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		ds.Database.Save(entity.Account{
			Id:             111,
			Name:           "test-occupied",
			Email:          "test-occupied@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 111})

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		occupied, err := repository.IsEmailOccupied("test-occupied@test.fundwit.com")
		assert.Equal(t, nil, err)
		assert.Equal(t, true, occupied)
	})
}

func TestDatabaseAccountRepository_IsEmailOccupied(it *testing.T) {
	it.Run("should return false if email is not occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		occupied, err := repository.IsEmailOccupied(uuid.New().String() + "@test.fundwit.com")
		assert.Equal(t, nil, err)
		assert.Equal(t, false, occupied)
	})

	it.Run("should return true if email is occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		ds.Database.Save(entity.Account{
			Id:             111,
			Name:           "test-occupied",
			Email:          "test-occupied@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 111})

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		occupied, err := repository.IsEmailOccupied("test-occupied@test.fundwit.com")
		assert.Equal(t, nil, err)
		assert.Equal(t, true, occupied)
	})
}

func TestDatabaseAccountRepository_IsAccountNameOccupied(it *testing.T) {
	it.Run("should return false if account name is not occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		occupied, err := repository.IsAccountNameOccupied(uuid.New().String())
		assert.Equal(t, nil, err)
		assert.Equal(t, false, occupied)
	})

	it.Run("should return true if account name is occupied", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		ds.Database.Save(entity.Account{
			Id:             111,
			Name:           "test-occupied",
			Email:          "test-occupied@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 111})

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		occupied, err := repository.IsAccountNameOccupied("test-occupied")
		assert.Equal(t, nil, err)
		assert.Equal(t, true, occupied)
	})
}

func TestDatabaseAccountRepository_FindByName(it *testing.T) {
	it.Run("should return nil if account name is not found", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}
		accountName := uuid.New().String()

		account, err := repository.FindByName(accountName)
		assert.Equal(t, true, gorm.IsRecordNotFoundError(err))
		assert.Nil(t, account)
	})

	it.Run("should return account if account name is found", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		account := entity.Account{
			Id:             111,
			Name:           "test-findByName",
			Email:          "test-occupied@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		}
		ds.Database.Save(account)
		defer ds.Database.Delete(entity.Account{Id: 111})

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		found, err := repository.FindByName(account.Name)
		assert.Equal(t, nil, err)
		assert.Equal(t, account.Id, found.Id)
	})
}

func TestDatabaseAccountRepository_Save(it *testing.T) {
	it.Run("should save account successfully", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		// 查询不到
		occupied, err := repository.IsAccountNameOccupied("test-save")
		assert.Equal(t, nil, err)
		assert.Equal(t, false, occupied)

		err = repository.Save(&entity.Account{
			Id:             222,
			Name:           "test-save",
			Email:          "test-save@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 222})
		assert.Equal(t, nil, err)

		// 可查询到
		occupied, err = repository.IsAccountNameOccupied("test-save")
		assert.Equal(t, nil, err)
		assert.Equal(t, true, occupied)
	})

	it.Run("should save failed with duplicated name, or email", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}

		err := repository.Save(&entity.Account{
			Id:             200,
			Name:           "test-save",
			Email:          "test-save@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		defer ds.Database.Delete(entity.Account{Id: 200})
		assert.Equal(t, nil, err)

		err = repository.Save(&entity.Account{
			Id:             201,
			Name:           "test-save",
			Email:          "test-save1@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		assert.Equal(t, "Error 1062: Duplicate entry 'test-save' for key 'name'", fmt.Sprintf("%s", err))

		err = repository.Save(&entity.Account{
			Id:             201,
			Name:           "test-save1",
			Email:          "test-save@test.fundwit.com",
			CreateTime:     time.Now(),
			LastUpdateTime: time.Now(),
		})
		assert.Equal(t, "Error 1062: Duplicate entry 'test-save@test.fundwit.com' for key 'email'", fmt.Sprintf("%s", err))
	})

	it.Run("should save failed when validate not pass", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		repository := &DatabaseAccountRepository{
			IdWorker: util.DefaultIdWorker,
			Database: ds.Database,
		}
		repository.Database.Delete(entity.Account{Name: "test-validate"})

		err := repository.Save(&entity.Account{Name: "test-validate", Email: "test-validate@test.fundwit.com", CreateTime: time.Now(), LastUpdateTime: time.Now()})
		assert.Equal(t, "Key: 'Account.Id' Error:Field validation for 'Id' failed on the 'required' tag", fmt.Sprintf("%s", err))

		err = repository.Save(&entity.Account{Id: util.DefaultIdWorker.NextIdOrFail(), Email: "test-validate@test.fundwit.com", CreateTime: time.Now(), LastUpdateTime: time.Now()})
		assert.Equal(t, "Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag", fmt.Sprintf("%s", err))

		err = repository.Save(&entity.Account{Id: util.DefaultIdWorker.NextIdOrFail(), Name: "test-validate", CreateTime: time.Now(), LastUpdateTime: time.Now()})
		assert.Equal(t, "Key: 'Account.Email' Error:Field validation for 'Email' failed on the 'required' tag", fmt.Sprintf("%s", err))

		err = repository.Save(&entity.Account{Id: util.DefaultIdWorker.NextIdOrFail(), Name: "test-validate", Email: "test-validate", CreateTime: time.Now(), LastUpdateTime: time.Now()})
		assert.Equal(t, "Key: 'Account.Email' Error:Field validation for 'Email' failed on the 'email' tag", fmt.Sprintf("%s", err))

		err = repository.Save(&entity.Account{Id: util.DefaultIdWorker.NextIdOrFail(), Name: "test-validate", Email: "test-validate@test.fundwit.com", LastUpdateTime: time.Now()})
		assert.Equal(t, "Key: 'Account.CreateTime' Error:Field validation for 'CreateTime' failed on the 'required' tag", fmt.Sprintf("%s", err))

		err = repository.Save(&entity.Account{Id: util.DefaultIdWorker.NextIdOrFail(), Name: "test-validate", Email: "test-validate@test.fundwit.com", CreateTime: time.Now()})
		assert.Equal(t, "Key: 'Account.LastUpdateTime' Error:Field validation for 'LastUpdateTime' failed on the 'required' tag", fmt.Sprintf("%s", err))
	})
}
