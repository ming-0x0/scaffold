package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/scaffold/pkg/domainerror"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// PartnersTableName TableName
var PartnersTableName = "partners"

type PartnerStatus int64

const (
	PartnerStatusActive   PartnerStatus = 1
	PartnerStatusInactive PartnerStatus = 2
)

// Partner struct
type Partner struct {
	ID            int64                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Name          string                      `gorm:"column:name;type:text;not null" mapstructure:"name"`
	DescriptionVi undefined.Undefined[string] `gorm:"column:description_vi;type:text" mapstructure:"description_vi"`
	DescriptionEn undefined.Undefined[string] `gorm:"column:description_en;type:text" mapstructure:"description_en"`
	DescriptionZh undefined.Undefined[string] `gorm:"column:description_zh;type:text" mapstructure:"description_zh"`
	Status        PartnerStatus               `gorm:"column:status;type:int;not null;default:1" mapstructure:"status"` // 1: active, 2: inactive
	Position      undefined.Undefined[int64]  `gorm:"column:position;type:int;unique" mapstructure:"position"`
	ResourceID    int64                       `gorm:"column:resource_id;type:bigint;not null" mapstructure:"resource_id"`
	Link          undefined.Undefined[string] `gorm:"column:link;type:text" mapstructure:"link"`
	BaseDomain
}

// TableName func
func (p *Partner) TableName() string {
	return PartnersTableName
}

func (p *Partner) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Status, validation.Required, validation.In(PartnerStatusActive, PartnerStatusInactive)),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type PartnerRepositoryInterface interface {
	repository.RepositoryInterface[*Partner]
}
