package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"hallo/domain/entity"
	"hallo/util"
)

const AccountTableName = "accounts"

//go:generate mockgen -destination AccountRepository_mock.go -package domain hallo/domain AccountRepository
type AccountRepository interface {
	NextId() (uint64, error)
	IsAccountNameOccupied(accountName string) (bool, error)
	IsEmailOccupied(accountName string) (bool, error)
	FindByName(accountName string) (*entity.Account, error)
	Count() (uint64, error)
	Save(account *entity.Account) error
}

type DatabaseAccountRepository struct {
	IdWorker *util.IdWorker
	Database *gorm.DB
}

func (repository *DatabaseAccountRepository) NextId() (uint64, error) {
	return repository.IdWorker.NextId()
}

func (repository *DatabaseAccountRepository) Count() (uint64, error) {
	var count uint64
	err := repository.Database.Table(AccountTableName).Count(&count).Error
	return count, err
}

// return (nil, gorm.ErrRecordNotFound) when account name is not found
func (repository *DatabaseAccountRepository) FindByName(accountName string) (*entity.Account, error) {
	account := &entity.Account{}
	if err := repository.Database.Table(AccountTableName).First(account, entity.Account{Name: accountName}).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (repository *DatabaseAccountRepository) IsAccountNameOccupied(accountName string) (bool, error) {
	rows, err := repository.Database.Table(AccountTableName).Select("id").Where(entity.Account{Name: accountName}).Limit(1).Rows()
	if err != nil {
		return true, errors.New("failed to query")
	}
	if rows.Next() {
		return true, nil
	}
	if rows.Err() != nil {
		return true, rows.Err()
	}
	return false, nil
}

func (repository *DatabaseAccountRepository) IsEmailOccupied(email string) (bool, error) {
	rows, err := repository.Database.Table(AccountTableName).Select("id").Where(entity.Account{Email: email}).Limit(1).Rows()
	if err != nil {
		return true, err
	}
	if rows.Next() {
		return true, nil
	}
	if rows.Err() != nil {
		return true, rows.Err()
	}
	return false, nil
}

func (repository *DatabaseAccountRepository) Save(account *entity.Account) error {
	validate := validator.New()
	if err := validate.Struct(account); err != nil {
		return err
	}
	return repository.Database.Save(account).Error
}
