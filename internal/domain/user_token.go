package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
)

// UserTokensTableName TableName
var UserTokensTableName = "user_tokens"

// UserToken struct
type UserToken struct {
	ID        int64     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null" mapstructure:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Token     string    `gorm:"column:token;type:text;not null" mapstructure:"token"`
	TokenID   string    `gorm:"column:token_id;type:text;not null" mapstructure:"token_id"`
	ExpiredAt time.Time `gorm:"column:expired_at;type:timestamp;not null" mapstructure:"expired_at"`
	BaseDomain
}

// TableName func
func (ut *UserToken) TableName() string {
	return UserTokensTableName
}

func (ut *UserToken) Validate() error {
	err := validation.ValidateStruct(ut,
		validation.Field(&ut.UserID, validation.Required),
		validation.Field(&ut.Token, validation.Required),
		validation.Field(&ut.TokenID, validation.Required),
		validation.Field(&ut.ExpiredAt, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type UserTokenRepositoryInterface interface {
	repository.RepositoryInterface[*UserToken]
}
