package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	*repository.Repository[*domain.User]
}

func NewUserRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *UserRepository {
	return &UserRepository{
		repository.New[*domain.User](db, logger),
	}
}
