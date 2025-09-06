package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FooterRepository struct {
	*repository.Repository[*domain.Footer]
}

func NewFooterRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *FooterRepository {
	return &FooterRepository{
		repository.New[*domain.Footer](db, logger),
	}
}
