package entity

import "time"

type Account struct {
	Id    uint64 `json:"id"    validate:"required"         gorm:"type:bigint;primary_key"                     pact:"example=10"`
	Name  string `json:"name"  validate:"required"         gorm:"type:nvarchar(127);unique;not null"          pact:"example=Sally"`
	Email string `json:"email" validate:"required,email"   gorm:"type:varchar(127);unique;not null"                    pact:"example=ann@test.com"`

	CreateTime     time.Time `json:"createTime"     validate:"required"    gorm:"type:DATETIME;not null"`
	LastUpdateTime time.Time `json:"lastUpdateTime" validate:"required"    grom:"type:DATETIME;not null"`
}

type EmailAccountCreateRequest struct {
	Email  string `json:"email"   validate:"required,email"  pact:"example=ann@test.com"`
	Name   string `json:"name"    validate:"required"        pact:"example=Sally"`
	Secret string `json:"secret"  validate:"required"   binding:"required"`
}
