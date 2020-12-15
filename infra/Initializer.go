package infra

import (
	"github.com/jinzhu/gorm"
	"hallo/domain/entity"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Account{})
	db.AutoMigrate(&entity.InternalIdentity{})
	db.AutoMigrate(&entity.IdentityBinding{})
}
