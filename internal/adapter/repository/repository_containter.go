package repository

import (
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	db     *gorm.DB
	logger *logrus.Logger
}

var _ domain.RepositoryContainerInterface = (*RepositoryContainer)(nil)

func NewRepositoryContainer(
	db *gorm.DB,
	logger *logrus.Logger,
) *RepositoryContainer {
	return &RepositoryContainer{
		db:     db,
		logger: logger,
	}
}

func (rc *RepositoryContainer) BannerRepository() domain.BannerRepositoryInterface {
	return NewBannerRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) CategoryRepository() domain.CategoryRepositoryInterface {
	return NewCategoryRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) CustomerRepository() domain.CustomerRepositoryInterface {
	return NewCustomerRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) FooterRepository() domain.FooterRepositoryInterface {
	return NewFooterRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) PartnerRepository() domain.PartnerRepositoryInterface {
	return NewPartnerRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) PermissionGroupRepository() domain.PermissionGroupRepositoryInterface {
	return NewPermissionGroupRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) PermissionRepository() domain.PermissionRepositoryInterface {
	return NewPermissionRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) PostRepository() domain.PostRepositoryInterface {
	return NewPostRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) ResourceRepository() domain.ResourceRepositoryInterface {
	return NewResourceRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) RoleRepository() domain.RoleRepositoryInterface {
	return NewRoleRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) UserRepository() domain.UserRepositoryInterface {
	return NewUserRepository(rc.db, rc.logger)
}

func (rc *RepositoryContainer) UserTokenRepository() domain.UserTokenRepositoryInterface {
	return NewUserTokenRepository(rc.db, rc.logger)
}
