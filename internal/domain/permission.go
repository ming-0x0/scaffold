package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
)

// PermissionsTableName TableName
var PermissionsTableName = "permissions"

// Permission struct
type Permission struct {
	ID             int64  `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	PermissionName string `gorm:"column:permission_name;type:varchar(255);not null;unique" mapstructure:"permission_name"`
	FunctionCode   string `gorm:"column:function_code;type:varchar(255);not null" mapstructure:"function_code"`
	BaseDomain
}

// TableName func
func (p *Permission) TableName() string {
	return PermissionsTableName
}

func (p *Permission) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.PermissionName, validation.Required),
		validation.Field(&p.FunctionCode, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type PermissionRepositoryInterface interface {
	repository.RepositoryInterface[*Permission]
}
