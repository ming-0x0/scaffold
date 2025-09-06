package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PermissionGroupRepository struct {
	*repository.Repository[*domain.PermissionGroup]
}

func NewPermissionGroupRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *PermissionGroupRepository {
	return &PermissionGroupRepository{
		repository.New[*domain.PermissionGroup](db, logger),
	}
}
