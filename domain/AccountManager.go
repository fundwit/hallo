package domain

import (
	"fmt"
	"hallo/domain/entity"
	"log"
	"time"
)

const InternalProviderId = "internal"

//go:generate mockgen -destination AccountManager_mock.go -package domain hallo/domain AccountManager
type AccountManager interface {
	CreateAccount(action entity.EmailAccountCreateRequest) (*entity.Account, error)
	AuthenticateInternalIdentity(accountName, secret string) (*entity.Account, error)
}

type AccountManagerImpl struct {
	AccountRepository          AccountRepository
	IdentityBindingRepository  IdentityBindingRepository
	InternalIdentityRepository InternalIdentityRepository
}

func (manager *AccountManagerImpl) CreateAccount(action entity.EmailAccountCreateRequest) (*entity.Account, error) {
	// TODO transaction
	isNameOccupied, err := manager.AccountRepository.IsAccountNameOccupied(action.Name)
	if err != nil {
		return nil, err
	}
	if isNameOccupied {
		return nil, &AccountNameIsOccupied{}
	}

	isEmailOccupied, err := manager.AccountRepository.IsEmailOccupied(action.Email)
	if err != nil {
		return nil, err
	}
	if isEmailOccupied {
		return nil, &AccountEmailIsOccupied{}
	}

	accountId, err := manager.AccountRepository.NextId()
	if err != nil {
		log.Println(err)
		return nil, IdGenerateFailure
	}

	now := time.Now()
	account := &entity.Account{
		Id:    accountId,
		Name:  action.Name,
		Email: action.Email,

		CreateTime:     now,
		LastUpdateTime: now,
	}

	err = manager.AccountRepository.Save(account)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// create binding (create internal identity for internalProvider)
	err = manager.bindIdentity(accountId, InternalProviderId, fmt.Sprintf("%d", accountId), action.Secret)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (manager *AccountManagerImpl) AuthenticateInternalIdentity(accountName, secret string) (*entity.Account, error) {
	// TODO tx begin
	account, err := manager.AccountRepository.FindByName(accountName)
	if err != nil {
		return nil, err
	}
	err = manager.InternalIdentityRepository.Authenticate(account.Id, secret)
	if err != nil {
		return nil, err
	}
	// TODO tx commit or rollback
	return account, nil
}

func (manager *AccountManagerImpl) bindIdentity(accountId uint64, providerId, providerAccountId, credential string) error {
	if providerId == InternalProviderId {
		// accountId and providerAccountId are equals, but in different type
		if err := manager.InternalIdentityRepository.Save(accountId, credential); err != nil {
			return err
		}
	}

	return manager.IdentityBindingRepository.Save(accountId, providerId, providerAccountId)
}
