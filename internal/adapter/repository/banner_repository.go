package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BannerRepository struct {
	*repository.Repository[*domain.Banner]
}

func NewBannerRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *BannerRepository {
	return &BannerRepository{
		repository.New[*domain.Banner](db, logger),
	}
}
