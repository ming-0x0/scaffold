package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserTokenRepository struct {
	*repository.Repository[*domain.UserToken]
}

func NewUserTokenRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *UserTokenRepository {
	return &UserTokenRepository{
		repository.New[*domain.UserToken](db, logger),
	}
}
