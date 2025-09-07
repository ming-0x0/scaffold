package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
)

// UsersTableName TableName
var UsersTableName = "users"

type UserStatus int32

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusInactive
)

type User struct {
	ID                int32           `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	FullName          string          `gorm:"column:full_name;type:varchar(50);not null" mapstructure:"full_name"`
	Email             string          `gorm:"column:email;type:varchar(255);not null;unique" mapstructure:"email"`
	Username          string          `gorm:"column:username;type:varchar(50);not null;unique" mapstructure:"username"`
	Password          string          `gorm:"column:password;type:text;not null" mapstructure:"password"`
	IsAdmin           bool            `gorm:"column:is_admin;type:tinyint(1);not null;default:0" mapstructure:"is_admin"`
	Status            UserStatus      `gorm:"column:status;type:integer;not null;default:1" mapstructure:"status"`                  // 1: active, 2: inactive
	ReceiveEmail      bool            `gorm:"column:receive_email;type:tinyint(1);not null;default:0" mapstructure:"receive_email"` // 0: no, 1: yes
	PermissionGroupID int32           `gorm:"column:permission_group_id;type:bigint;not null" mapstructure:"permission_group_id"`
	PermissionGroup   PermissionGroup `gorm:"foreignKey:PermissionGroupID"`
	BaseDomain
}

// TableName func
func (u *User) TableName() string {
	return UsersTableName
}

func (u *User) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.FullName, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.IsAdmin, validation.In(true, false)),
		validation.Field(&u.Status, validation.In(UserStatusActive, UserStatusInactive)),
		validation.Field(&u.ReceiveEmail, validation.In(true, false)),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type UserRepositoryInterface interface {
	repository.RepositoryInterface[*User]
}
