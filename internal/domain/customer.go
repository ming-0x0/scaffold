package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/ming-0x0/scaffold/pkg/undefined"
)

// CustomersTableName TableName
var CustomersTableName = "customers"

type CustomerStatus int32

const (
	CustomerStatusActive   CustomerStatus = 1 // đã trả lời
	CustomerStatusInactive CustomerStatus = 2 // chưa trả lời
)

type ServiceType int32

const (
	ServiceTypeTuyenDung ServiceType = 1
	ServiceTypeLienHe    ServiceType = 2
	ServiceTypeKhoaHoc   ServiceType = 3
)

// Customer struct
type Customer struct {
	ID           int32                       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	CustomerName string                      `gorm:"column:customer_name;type:text;not null" mapstructure:"customer_name"`
	Email        string                      `gorm:"column:email;type:text;not null" mapstructure:"email"`
	PhoneNumber  string                      `gorm:"column:phone_number;type:text;not null" mapstructure:"phone_number"`
	CompanyName  undefined.Undefined[string] `gorm:"column:company_name;type:text" mapstructure:"company_name"`
	Message      undefined.Undefined[string] `gorm:"column:message;type:text" mapstructure:"message"`
	Note         undefined.Undefined[string] `gorm:"column:note;type:text" mapstructure:"note"`
	ServiceType  ServiceType                 `gorm:"column:service_type;type:int;not null" mapstructure:"service_type"`
	Status       CustomerStatus              `gorm:"column:status;type:int;not null;default:2" mapstructure:"status"`
	BaseDomain
}

// TableName TableName
func (c *Customer) TableName() string {
	return CustomersTableName
}

func (c *Customer) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.CustomerName, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.PhoneNumber, validation.Required),
		validation.Field(&c.ServiceType, validation.Required, validation.In(ServiceTypeTuyenDung, ServiceTypeLienHe, ServiceTypeKhoaHoc)),
		validation.Field(&c.Status, validation.Required, validation.In(CustomerStatusActive, CustomerStatusInactive)),
	)
	if err != nil {
		return domainerror.Wrap(domainerror.InvalidArgument, err)
	}

	return nil
}

type CustomerRepositoryInterface interface {
	repository.RepositoryInterface[*Customer]
}
