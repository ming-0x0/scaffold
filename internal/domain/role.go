package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
)

// RolesTableName TableName
var RolesTableName = "roles"

// Role struct
type Role struct {
	ID                int64           `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	PermissionID      int64           `gorm:"column:permission_id;type:bigint;not null" mapstructure:"permission_id"`
	Permission        Permission      `gorm:"foreignKey:PermissionID"`
	PermissionGroupID int64           `gorm:"column:permission_group_id;type:bigint;not null" mapstructure:"permission_group_id"`
	PermissionGroup   PermissionGroup `gorm:"foreignKey:PermissionGroupID"`
	BaseDomain
}

// TableName func
func (r *Role) TableName() string {
	return RolesTableName
}

func (r *Role) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.PermissionID, validation.Required),
		validation.Field(&r.PermissionGroupID, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type RoleRepositoryInterface interface {
	repository.RepositoryInterface[*Role]
}
