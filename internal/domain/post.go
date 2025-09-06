package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/pkg/domainerror"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// PostsTableName TableName
var PostsTableName = "posts"

type PostStatus int64

const (
	PostStatusDraft PostStatus = iota + 1
	PostStatusInReview
	PostStatusPublic
	PostStatusReject
	PostStatusRemoved
)

type PostType int64

const (
	PostTypeNews PostType = iota + 1
	PostTypeProduct
	PostTypeProject
	PostTypeTechnology
	PostTypeDelivery
)

// Post struct
type Post struct {
	ID            int64                          `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	TitleVi       string                         `gorm:"column:title_vi;type:text;not null" mapstructure:"title_vi"`
	TitleEn       string                         `gorm:"column:title_en;type:text;not null" mapstructure:"title_en"`
	TitleZh       string                         `gorm:"column:title_zh;type:text;not null" mapstructure:"title_zh"`
	SlugVi        string                         `gorm:"column:slug_vi;type:text;not null" mapstructure:"slug_vi"`
	SlugEn        string                         `gorm:"column:slug_en;type:text;not null" mapstructure:"slug_en"`
	SlugZh        string                         `gorm:"column:slug_zh;type:text;not null" mapstructure:"slug_zh"`
	AltVi         string                         `gorm:"column:alt_vi;type:text;not null" mapstructure:"alt_vi"`
	AltEn         string                         `gorm:"column:alt_en;type:text;not null" mapstructure:"alt_en"`
	AltZh         string                         `gorm:"column:alt_zh;type:text;not null" mapstructure:"alt_zh"`
	DescriptionVi undefined.Undefined[string]    `gorm:"column:description_vi;type:text" mapstructure:"description_vi"`
	DescriptionEn undefined.Undefined[string]    `gorm:"column:description_en;type:text" mapstructure:"description_en"`
	DescriptionZh undefined.Undefined[string]    `gorm:"column:description_zh;type:text" mapstructure:"description_zh"`
	Avatar        int64                          `gorm:"column:avatar;type:bigint;not null" mapstructure:"avatar"` // ảnh đại diện
	ResourceIDs   undefined.Undefined[string]    `gorm:"column:resource_ids;type:text" mapstructure:"resource_ids"`
	ContentVi     string                         `gorm:"column:content_vi;type:text;not null" mapstructure:"content_vi"`
	ContentEn     string                         `gorm:"column:content_en;type:text;not null" mapstructure:"content_en"`
	ContentZh     string                         `gorm:"column:content_zh;type:text;not null" mapstructure:"content_zh"`
	Status        PostStatus                     `gorm:"column:status;type:int;not null" mapstructure:"status"`                      // 1: lưu nháp, 2: chờ duyệt, 3: hoạt động, 4: từ chối, 5: gỡ bài
	Type          PostType                       `gorm:"column:type;type:int;not null;default:1" mapstructure:"type"`                // 1: Tin tức, 2: Sản phẩm, 3: Dự án
	Flagship      bool                           `gorm:"column:flagship;type:tinyint(1);not null;default:0" mapstructure:"flagship"` // nổi bật
	ColorPalette  undefined.Undefined[string]    `gorm:"column:color_palette;type:text" mapstructure:"color_palette"`                // bảng màu
	CategoryID    int                            `gorm:"column:category_id;type:bigint;not null" mapstructure:"category_id"`
	Category      Category                       `gorm:"foreignKey:CategoryID;references:ID"`
	PublicDate    undefined.Undefined[time.Time] `gorm:"column:public_date;type:timestamp" mapstructure:"public_date"`
	InfoVi        undefined.Undefined[string]    `gorm:"column:info_vi;type:text" mapstructure:"info_vi"`
	InfoEn        undefined.Undefined[string]    `gorm:"column:info_en;type:text" mapstructure:"info_en"`
	InfoZh        undefined.Undefined[string]    `gorm:"column:info_zh;type:text" mapstructure:"info_zh"`
	BaseDomain
}

// TableName func
func (p *Post) TableName() string {
	return PostsTableName
}

func (p *Post) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.TitleVi, validation.Required),
		validation.Field(&p.TitleEn, validation.Required),
		validation.Field(&p.TitleZh, validation.Required),
		validation.Field(&p.SlugVi, validation.Required),
		validation.Field(&p.SlugEn, validation.Required),
		validation.Field(&p.SlugZh, validation.Required),
		validation.Field(&p.AltVi, validation.Required),
		validation.Field(&p.AltEn, validation.Required),
		validation.Field(&p.AltZh, validation.Required),
		validation.Field(&p.Avatar, validation.Required),
		validation.Field(&p.ContentVi, validation.Required),
		validation.Field(&p.ContentEn, validation.Required),
		validation.Field(&p.ContentZh, validation.Required),
		validation.Field(&p.Status, validation.Required, validation.In(PostStatusDraft, PostStatusInReview, PostStatusPublic, PostStatusReject, PostStatusRemoved)),
		validation.Field(&p.Type, validation.Required, validation.In(PostTypeNews, PostTypeProduct, PostTypeProject, PostTypeTechnology, PostTypeDelivery)),
		validation.Field(&p.Flagship, validation.Required),
		validation.Field(&p.CategoryID, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type PostRepositoryInterface interface {
	repository.RepositoryInterface[*Post]
}
