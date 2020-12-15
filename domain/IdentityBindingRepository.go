package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"hallo/domain/entity"
	"time"
)

const IdentityBindingTableName = "identity_bindings"

//go:generate mockgen -destination IdentityBindingRepository_mock.go -package domain hallo/domain IdentityBindingRepository
type IdentityBindingRepository interface {
	Save(accountId uint64, providerId, providerAccountId string) error
}

type DatabaseIdentityBindingRepository struct {
	Database *gorm.DB
}

func (repository *DatabaseIdentityBindingRepository) Save(accountId uint64, providerId, providerAccountId string) error {
	identityBinding := entity.IdentityBinding{
		AccountId:         accountId,
		ProviderId:        providerId,
		ProviderAccountId: providerAccountId,
		CreateTime:        time.Now(),
	}

	validate := validator.New()
	if err := validate.Struct(identityBinding); err != nil {
		return err
	}

	return repository.Database.Save(identityBinding).Error
}
