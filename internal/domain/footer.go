package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// FootersTableName TableName
var FootersTableName = "footers"

// Footer struct
type Footer struct {
	ID        int64                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	NameVi    string                      `gorm:"column:name_vi;type:text;not null" mapstructure:"name_vi"`
	NameEn    string                      `gorm:"column:name_en;type:text;not null" mapstructure:"name_en"`
	NameZh    string                      `gorm:"column:name_zh;type:text;not null" mapstructure:"name_zh"`
	ContentVi undefined.Undefined[string] `gorm:"column:content_vi;type:text" mapstructure:"content_vi"`
	ContentEn undefined.Undefined[string] `gorm:"column:content_en;type:text" mapstructure:"content_en"`
	ContentZh undefined.Undefined[string] `gorm:"column:content_zh;type:text" mapstructure:"content_zh"`
	Link      undefined.Undefined[string] `gorm:"column:link;type:text" mapstructure:"link"`
	BaseDomain
}

// TableName func
func (f *Footer) TableName() string {
	return FootersTableName
}

func (f *Footer) Validate() error {
	err := validation.ValidateStruct(f,
		validation.Field(&f.NameVi, validation.Required),
		validation.Field(&f.NameEn, validation.Required),
		validation.Field(&f.NameZh, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type FooterRepositoryInterface interface {
	repository.RepositoryInterface[*Footer]
}
