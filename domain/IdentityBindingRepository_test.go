package domain

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"hallo/domain/entity"
	"hallo/testinfra"
	"hallo/util"
	"testing"
	"time"
)

func TestDatabaseIdentityBindingRepository_Save(it *testing.T) {
	it.Run("should save identity binding successfully", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		store := &DatabaseIdentityBindingRepository{
			Database: ds.Database,
		}

		// 查询不到
		binding := entity.IdentityBinding{AccountId: util.DefaultIdWorker.NextIdOrFail(), ProviderId: "123", ProviderAccountId: "123"}
		err := store.Save(binding.AccountId, binding.ProviderId, binding.ProviderAccountId)
		assert.Equal(t, nil, err)

		// 验证绑定已存在
		rows, err := store.Database.Table(IdentityBindingTableName).Select("1").Where(entity.IdentityBinding{AccountId: binding.AccountId}).Limit(1).Rows()
		if err != nil || !rows.Next() {
			t.Fatal("save failed")
		}
	})

	it.Run("should save failed with database integrity issues", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		store := &DatabaseIdentityBindingRepository{
			Database: ds.Database,
		}

		// 插入数据
		providerId := "database-integrity"
		binding := entity.IdentityBinding{AccountId: util.DefaultIdWorker.NextIdOrFail(), ProviderId: providerId, ProviderAccountId: "123", CreateTime: time.Now()}
		err := store.Save(binding.AccountId, binding.ProviderId, binding.ProviderAccountId)
		assert.Equal(t, nil, err)

		// 插入复合主键重复数据
		binding = entity.IdentityBinding{AccountId: binding.AccountId, ProviderId: providerId, ProviderAccountId: "123", CreateTime: time.Now()}
		err = store.Database.Create(binding).Error

		assert.Equal(t, "Error 1062: Duplicate entry '123-database-integrity-"+fmt.Sprintf("%d", binding.AccountId)+"' for key 'PRIMARY'", fmt.Sprintf("%s", err))
	})

	it.Run("should save failed when validate not pass", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		store := &DatabaseIdentityBindingRepository{
			Database: ds.Database,
		}
		provider := uuid.New().String()
		store.Database.Delete(entity.IdentityBinding{ProviderId: provider})

		binding := entity.IdentityBinding{ProviderId: provider, ProviderAccountId: "123"}
		err := store.Save(binding.AccountId, binding.ProviderId, binding.ProviderAccountId)
		assert.Equal(t, "Key: 'IdentityBinding.AccountId' Error:Field validation for 'AccountId' failed on the 'required' tag", fmt.Sprintf("%s", err))

		binding = entity.IdentityBinding{AccountId: 123, ProviderAccountId: "123"}
		err = store.Save(binding.AccountId, binding.ProviderId, binding.ProviderAccountId)
		assert.Equal(t, "Key: 'IdentityBinding.ProviderId' Error:Field validation for 'ProviderId' failed on the 'required' tag", fmt.Sprintf("%s", err))

		binding = entity.IdentityBinding{AccountId: 123, ProviderId: provider}
		err = store.Save(binding.AccountId, binding.ProviderId, binding.ProviderAccountId)
		assert.Equal(t, "Key: 'IdentityBinding.ProviderAccountId' Error:Field validation for 'ProviderAccountId' failed on the 'required' tag", fmt.Sprintf("%s", err))
	})
}
