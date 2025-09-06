package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	*repository.Repository[*domain.Permission]
}

func NewPermissionRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *PermissionRepository {
	return &PermissionRepository{
		repository.New[*domain.Permission](db, logger),
	}
}
