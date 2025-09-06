package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PartnerRepository struct {
	*repository.Repository[*domain.Partner]
}

func NewPartnerRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *PartnerRepository {
	return &PartnerRepository{
		repository.New[*domain.Partner](db, logger),
	}
}
