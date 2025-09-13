package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// BannersTableName TableName
var BannersTableName = "banners"

type BannerStatus int32

const (
	BannerStatusActive   BannerStatus = 1
	BannerStatusInactive BannerStatus = 2
)

type BannerType int32

const (
	NoRefBannerType      BannerType = 1
	PostRefBannerType    BannerType = 2
	ServiceRefBannerType BannerType = 3
	ContactRefBannerType BannerType = 4
	CourseRefBannerType  BannerType = 5
)

//go:generate go run ../../cmd/fieldgen/main.go -struct=Banner -input=banner.go -output=../adapter/banner/banner_columns.go -table=banners -package=columns
type Banner struct {
	ID            int32                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	NameVi        string                      `gorm:"column:name_vi;type:text;not null" mapstructure:"name_vi"`
	NameEn        string                      `gorm:"column:name_en;type:text;not null" mapstructure:"name_en"`
	NameZh        string                      `gorm:"column:name_zh;type:text;not null" mapstructure:"name_zh"`
	DescriptionVi undefined.Undefined[string] `gorm:"column:description_vi;type:text" mapstructure:"description_vi"`
	DescriptionEn undefined.Undefined[string] `gorm:"column:description_en;type:text" mapstructure:"description_en"`
	DescriptionZh undefined.Undefined[string] `gorm:"column:description_zh;type:text" mapstructure:"description_zh"`
	Position      undefined.Undefined[int32]  `gorm:"column:position;type:int;unique" mapstructure:"position"`
	Status        BannerStatus                `gorm:"column:status;type:int;not null;default:1" mapstructure:"status"`
	ResourceID    int32                       `gorm:"column:resource_id;type:bigint;not null" mapstructure:"resource_id"`
	Link          undefined.Undefined[string] `gorm:"column:link;type:text" mapstructure:"link"`
	ButtonNameVi  undefined.Undefined[string] `gorm:"column:button_name_vi;type:text" mapstructure:"button_name_vi"`
	ButtonNameEn  undefined.Undefined[string] `gorm:"column:button_name_en;type:text" mapstructure:"button_name_en"`
	ButtonNameZh  undefined.Undefined[string] `gorm:"column:button_name_zh;type:text" mapstructure:"button_name_zh"`
	HasContent    bool                        `gorm:"column:has_content;type:tinyint(1);not null;default:0" mapstructure:"has_content"`
	BaseDomain
}

func (b *Banner) TableName() string {
	return BannersTableName
}

func (b *Banner) Validate() error {
	err := validation.ValidateStruct(b,
		validation.Field(&b.Status, validation.Required, validation.In(BannerStatusActive, BannerStatusInactive)),
		validation.Field(&b.ResourceID, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type BannerRepositoryInterface interface {
	repository.RepositoryInterface[*Banner]
}
