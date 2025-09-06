package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/ming-0x0/scaffold/pkg/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostRepository struct {
	*repository.Repository[*domain.Post]
}

func NewPostRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) *PostRepository {
	return &PostRepository{
		repository.New[*domain.Post](db, logger),
	}
}
