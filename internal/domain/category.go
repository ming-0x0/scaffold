package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// CategoriesTableName TableName
var CategoriesTableName = "categories"

type CategoryType int64

const (
	CategoryTypeDetail CategoryType = iota + 2
	CategoryTypeList
)

type CategoryStatus int64

const (
	CategoryStatusActive CategoryStatus = iota + 1
	CategoryStatusInactive
)

type Category struct {
	ID              int64                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	NameVi          string                      `gorm:"column:name_vi;type:text;not null" mapstructure:"name_vi"`
	NameEn          string                      `gorm:"column:name_en;type:text;not null" mapstructure:"name_en"`
	NameZh          string                      `gorm:"column:name_zh;type:text;not null" mapstructure:"name_zh"`
	DescriptionVi   undefined.Undefined[string] `gorm:"column:description_vi;type:text" mapstructure:"description_vi"`
	DescriptionEn   undefined.Undefined[string] `gorm:"column:description_en;type:text" mapstructure:"description_en"`
	DescriptionZh   undefined.Undefined[string] `gorm:"column:description_zh;type:text" mapstructure:"description_zh"`
	Type            CategoryType                `gorm:"column:type;type:bigint;not null;idx_parent_type_position " mapstructure:"type"`
	Editable        bool                        `gorm:"column:editable;type:tinyint(1);not null;default:1" mapstructure:"editable"`
	RouterVi        undefined.Undefined[string] `gorm:"column:router_vi;type:text" mapstructure:"router_vi"`
	RouterEn        undefined.Undefined[string] `gorm:"column:router_en;type:text" mapstructure:"router_en"`
	RouterZh        undefined.Undefined[string] `gorm:"column:router_zh;type:text" mapstructure:"router_zh"`
	Position        undefined.Undefined[int64]  `gorm:"column:position;type:bigint;idx_parent_type_position" mapstructure:"position"`
	ParentID        undefined.Undefined[int64]  `gorm:"column:parent_id;type:bigint;idx_parent_type_position" mapstructure:"parent_id"`
	Status          CategoryStatus              `gorm:"column:status;type:int;not null;default:1" mapstructure:"status"`
	ResourceID      undefined.Undefined[int64]  `gorm:"column:resource_id;type:bigint" mapstructure:"resource_id"`
	Level           int64                       `gorm:"column:level;type:bigint;not null" mapstructure:"level"`
	ChildCategories []Category                  `gorm:"foreignKey:ParentID"`
	BaseDomain
}

func (c *Category) TableName() string {
	return CategoriesTableName
}

func (c *Category) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.NameVi, validation.Required),
		validation.Field(&c.NameEn, validation.Required),
		validation.Field(&c.NameZh, validation.Required),
		validation.Field(&c.Type, validation.Required, validation.In(CategoryTypeDetail, CategoryTypeList)),
		validation.Field(&c.Level, validation.Required),
		validation.Field(&c.Status, validation.Required, validation.In(CategoryStatusActive, CategoryStatusInactive)),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type CategoryRepositoryInterface interface {
	repository.RepositoryInterface[*Category]
}
