package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"hallo/domain/entity"
	"hallo/testinfra"
	"hallo/util"
	"testing"
)

func TestDatabaseInternalIdentityRepository_hashCredential(t *testing.T) {
	t.Run("should hash credential successfully", func(t *testing.T) {
		// prepare test data
		store := &DatabaseInternalIdentityRepository{}
		assert.Equal(t, store.hashCredential("12345密码"), "b40a2ae84db3800da91c40ba920808cc4942b929")
	})
}

func TestDatabaseInternalIdentityRepository_Authenticate(t *testing.T) {
	//os.Setenv("DOCKER_HOST", "tcp://192.168.2.108:2375")
	//mysqlService, err := testinfra.NewMysqlContainer()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer mysqlService.Stop()
	//mysqlServer := fmt.Sprintf("%s:%s@(%s:%s)", mysqlService.Username, mysqlService.Password, mysqlService.Host, mysqlService.MappedPort.Port())

	t.Run("it should authentication failed when user identity is not exist", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		// prepare test data
		store := &DatabaseInternalIdentityRepository{
			Database: ds.Database,
		}
		accountId := uint64(123)
		credential := "123456"

		// pre-assertion: not exist
		rows, err := store.Database.Model(&entity.InternalIdentity{}).Select("1").
			Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(credential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), false)

		// do authentication
		err = store.Authenticate(accountId, credential)
		assert.Equal(t, err, &AccountAuthenticationFailure{})

		// do save
		err = store.Save(accountId, credential)
		assert.Equal(t, err, nil)

		// do authentication with invalid credential
		err = store.Authenticate(accountId, credential+"bad")
		assert.Equal(t, err, &AccountAuthenticationFailure{})

		// do authentication successfully
		err = store.Authenticate(accountId, credential)
		// do assertion
		assert.Nil(t, err)
	})
}

func TestDatabaseInternalIdentityRepository_Delete(t *testing.T) {
	t.Run("it should delete user identity successfully", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		// prepare test data
		realm := &DatabaseInternalIdentityRepository{
			Database: ds.Database,
		}
		accountId := uint64(123)
		credential := "123456"
		err := realm.Save(accountId, credential)

		// pre-assertion
		rows, err := realm.Database.Model(&entity.InternalIdentity{}).
			Select("1").Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(credential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), true)

		// do delete
		err = realm.Delete(accountId)
		// do assertion
		if err != nil {
			t.Errorf("unepxected error occured when do delete %v", err)
		}
		rows, err = realm.Database.Model(&entity.InternalIdentity{}).Select("1").
			Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(credential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), false)
	})
}

func TestDatabaseInternalIdentityRepository_Save(t *testing.T) {
	t.Run("it should save new user identity successfully", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		// prepare test data
		realm := &DatabaseInternalIdentityRepository{
			Database: ds.Database,
		}
		accountId := uint64(123)
		credential := "123456"

		rows, err := realm.Database.Model(&entity.InternalIdentity{}).Select("1").
			Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(credential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), false)

		// do insert
		err = realm.Save(accountId, credential)

		// do assertion
		if err != nil {
			t.Errorf("unepxected error occured when do save %v", err)
		}
		rows, err = realm.Database.Model(&entity.InternalIdentity{}).Select("1").
			Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(credential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), true)

		newCredential := "654321"
		// do update
		err = realm.Save(accountId, newCredential)
		if err != nil {
			t.Errorf("unepxected error occured when do save %v", err)
		}

		rows, err = realm.Database.Model(&entity.InternalIdentity{}).Select("1").
			Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(credential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), false)

		rows, err = realm.Database.Model(&entity.InternalIdentity{}).Select("1").
			Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: util.HashSha1Hex([]byte(newCredential))}).Limit(1).Rows()
		assert.Equal(t, err, nil)
		assert.Equal(t, rows.Next(), true)
	})

	t.Run("should save failed when validate not pass", func(t *testing.T) {
		ds := testinfra.NewTemporaryDatabase()
		defer ds.CleanAndDisconnect()

		store := &DatabaseInternalIdentityRepository{
			Database: ds.Database,
		}

		err := store.Save(0, "123")
		assert.Equal(t, "Key: 'InternalIdentity.AccountId' Error:Field validation for 'AccountId' failed on the 'required' tag", fmt.Sprintf("%s", err))

		err = store.Save(123, "")
		assert.Equal(t, "Key: '' Error:Field validation for '' failed on the 'required' tag", fmt.Sprintf("%s", err))
	})
}
