package entity

import "time"

type IdentityBinding struct {
	ProviderAccountId string `validate:"required" gorm:"type:nvarchar(127);primary_key"`
	ProviderId        string `validate:"required" gorm:"type:nvarchar(127);primary_key"`
	AccountId         uint64 `validate:"required" gorm:"type:bigint;primary_key"`

	CreateTime time.Time `validate:"required" grom:"type:DATETIME;not null"`
}
