package adapter

import (
	"github.com/ming-0x0/scaffold/internal/adapter/repository"
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Adapter struct {
	db     *gorm.DB
	logger *logrus.Logger
}

var _ domain.AdapterInterface = (*Adapter)(nil)

func New(
	db *gorm.DB,
	logger *logrus.Logger,
) *Adapter {
	return &Adapter{
		db:     db,
		logger: logger,
	}
}

func (a *Adapter) BannerRepository() domain.BannerRepositoryInterface {
	return repository.NewBannerRepository(a.db, a.logger)
}

func (a *Adapter) CategoryRepository() domain.CategoryRepositoryInterface {
	return repository.NewCategoryRepository(a.db, a.logger)
}

func (a *Adapter) CustomerRepository() domain.CustomerRepositoryInterface {
	return repository.NewCustomerRepository(a.db, a.logger)
}

func (a *Adapter) FooterRepository() domain.FooterRepositoryInterface {
	return repository.NewFooterRepository(a.db, a.logger)
}

func (a *Adapter) PartnerRepository() domain.PartnerRepositoryInterface {
	return repository.NewPartnerRepository(a.db, a.logger)
}

func (a *Adapter) PermissionGroupRepository() domain.PermissionGroupRepositoryInterface {
	return repository.NewPermissionGroupRepository(a.db, a.logger)
}

func (a *Adapter) PermissionRepository() domain.PermissionRepositoryInterface {
	return repository.NewPermissionRepository(a.db, a.logger)
}

func (a *Adapter) PostRepository() domain.PostRepositoryInterface {
	return repository.NewPostRepository(a.db, a.logger)
}

func (a *Adapter) ResourceRepository() domain.ResourceRepositoryInterface {
	return repository.NewResourceRepository(a.db, a.logger)
}

func (a *Adapter) RoleRepository() domain.RoleRepositoryInterface {
	return repository.NewRoleRepository(a.db, a.logger)
}

func (a *Adapter) UserRepository() domain.UserRepositoryInterface {
	return repository.NewUserRepository(a.db, a.logger)
}

func (a *Adapter) UserTokenRepository() domain.UserTokenRepositoryInterface {
	return repository.NewUserTokenRepository(a.db, a.logger)
}
