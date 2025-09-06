package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/pkg/domainerror"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// PermissionGroupsTableName TableName
var PermissionGroupsTableName = "permission_groups"

// PermissionGroup struct
type PermissionGroup struct {
	ID             int64                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Name           string                      `gorm:"column:name;type:varchar(255);not null;unique" mapstructure:"name"`
	Description    undefined.Undefined[string] `json:",omitzero"`
	FullPermission bool                        `gorm:"column:full_permission;type:tinyint(1);not null;default:0" mapstructure:"full_permission"`
	BaseDomain
}

// TableName func
func (pg *PermissionGroup) TableName() string {
	return PermissionGroupsTableName
}

func (pg *PermissionGroup) Validate() error {
	err := validation.ValidateStruct(pg,
		validation.Field(&pg.Name, validation.Required),
		validation.Field(&pg.FullPermission, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type PermissionGroupRepositoryInterface interface {
	repository.RepositoryInterface[*PermissionGroup]
}
