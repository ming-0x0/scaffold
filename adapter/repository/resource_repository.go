package repository

import (
	"github.com/ming-0x0/scaffold/domain"
	"github.com/ming-0x0/scaffold/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ResourceRepository struct {
	*repository.Repository[*domain.Resource]
}

func NewResourceRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *ResourceRepository {
	return &ResourceRepository{
		repository.New[*domain.Resource](db, logger),
	}
}
