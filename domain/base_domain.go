package domain

import (
	"time"

	"gorm.io/gorm"
)

type BaseDomain struct {
	CreatedAt time.Time `gorm:"column:created_at;not null;type:datetime;default:current_timestamp" mapstructure:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;type:datetime;default:current_timestamp" mapstructure:"updated_at" json:"updated_at"`
}

type BaseDomainWithDeleted struct {
	BaseDomain
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index" mapstructure:"deleted_at" json:"deleted_at"`
}
