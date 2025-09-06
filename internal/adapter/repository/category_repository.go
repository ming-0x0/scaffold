package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	*repository.Repository[*domain.Category]
}

func NewCategoryRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *CategoryRepository {
	return &CategoryRepository{
		repository.New[*domain.Category](db, logger),
	}
}
