package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository struct {
	*repository.Repository[*domain.Role]
}

func NewRoleRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *RoleRepository {
	return &RoleRepository{
		repository.New[*domain.Role](db, logger),
	}
}
