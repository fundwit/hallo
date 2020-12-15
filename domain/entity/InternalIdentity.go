package entity

import "time"

type InternalIdentity struct {
	AccountId      uint64    `validate:"required" gorm:"type:bigint;primary_key"`
	HashedIdentity string    `validate:"required" gorm:"type:nvarchar(255);not null"`
	CreateTime     time.Time `validate:"required" grom:"type:DATETIME;not null"`
}
