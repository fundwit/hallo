package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"hallo/domain/entity"
	"hallo/util"
	"log"
	"time"
)

//go:generate mockgen -destination InternalIdentityRepository_mock.go -package domain hallo/domain InternalIdentityRepository
type InternalIdentityRepository interface {
	Save(accountId uint64, credential string) error
	Authenticate(accountId uint64, credential string) error
}

type DatabaseInternalIdentityRepository struct {
	Database *gorm.DB
}

func (repository *DatabaseInternalIdentityRepository) Save(accountId uint64, credential string) error {
	validate := validator.New()
	if err := validate.Var(credential, "required"); err != nil {
		return err
	}

	internalIdentity := entity.InternalIdentity{
		AccountId:      accountId,
		HashedIdentity: repository.hashCredential(credential),
		CreateTime:     time.Now(),
	}

	if err := validate.Struct(internalIdentity); err != nil {
		return err
	}

	return repository.Database.Save(internalIdentity).Error
}

func (repository *DatabaseInternalIdentityRepository) Authenticate(accountId uint64, credential string) error {
	hashedIdentity := repository.hashCredential(credential)
	rows, err := repository.Database.Model(&entity.InternalIdentity{}).Select("1").
		Where(entity.InternalIdentity{AccountId: accountId, HashedIdentity: hashedIdentity}).Limit(1).Rows()
	if err != nil {
		return err
	}
	if rows.Next() {
		return nil
	} else {
		return &AccountAuthenticationFailure{}
	}
}

func (repository *DatabaseInternalIdentityRepository) Delete(accountId uint64) error {
	db := repository.Database.Where(entity.InternalIdentity{AccountId: accountId}).Delete(&entity.InternalIdentity{})
	log.Printf("delete identity of account [%d] -> success, rows affected [%d].\n", accountId, db.RowsAffected)
	return db.Error
}

func (repository *DatabaseInternalIdentityRepository) hashCredential(credential string) string {
	return util.HashSha1Hex([]byte(credential))
}
