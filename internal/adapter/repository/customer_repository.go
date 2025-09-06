package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	*repository.Repository[*domain.Customer]
}

func NewCustomerRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *CustomerRepository {
	return &CustomerRepository{
		repository.New[*domain.Customer](db, logger),
	}
}
