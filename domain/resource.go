package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/shared/domainerror"
	"github.com/ming-0x0/scaffold/shared/repository"
	"github.com/ming-0x0/scaffold/shared/undefined"
)

// ResourcesTableName TableName
var ResourcesTableName = "resources"

type ResourceType int64

const (
	ResourceTypeImage ResourceType = iota + 1
	ResourceTypeVideo
)

type Resource struct {
	ID          int64                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Name        string                      `gorm:"column:name;type:varchar(255);not null;unique" mapstructure:"name"`
	Description undefined.Undefined[string] `gorm:"column:description;type:longtext" mapstructure:"description"`
	Type        ResourceType                `gorm:"column:type;type:bigint;not null" mapstructure:"type"`
	Url         string                      `gorm:"column:url;type:longtext;not null" mapstructure:"url"`
	YoutubeID   undefined.Undefined[string] `gorm:"column:youtube_id;type:text" mapstructure:"youtube_id"`
	Banners     []Banner                    `gorm:"foreignKey:ResourceID;contact:onDelete:RESTRICT"`
	Categories  []Category                  `gorm:"foreignKey:ResourceID;contact:onDelete:RESTRICT"`
	Posts       []Post                      `gorm:"foreignKey:Avatar;contact:onDelete:RESTRICT"`
	BaseDomain
}

func (r *Resource) TableName() string {
	return ResourcesTableName
}

func (r *Resource) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Type, validation.Required, validation.In(ResourceTypeImage, ResourceTypeVideo)),
		validation.Field(&r.Url, validation.Required),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type ResourceRepositoryInterface interface {
	repository.RepositoryInterface[*Resource]
}
