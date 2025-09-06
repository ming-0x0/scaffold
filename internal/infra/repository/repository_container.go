package repository

import "github.com/ming-0x0/scaffold/internal/domain"

type RepositoryContainerInterface interface {
	BannerRepository() domain.BannerRepositoryInterface
	CategoryRepository() domain.CategoryRepositoryInterface
	CustomerRepository() domain.CustomerRepositoryInterface
	FooterRepository() domain.FooterRepositoryInterface
	PartnerRepository() domain.PartnerRepositoryInterface
	PermissionGroupRepository() domain.PermissionGroupRepositoryInterface
	PermissionRepository() domain.PermissionRepositoryInterface
	PostRepository() domain.PostRepositoryInterface
	ResourceRepository() domain.ResourceRepositoryInterface
	RoleRepository() domain.RoleRepositoryInterface
	UserRepository() domain.UserRepositoryInterface
	UserTokenRepository() domain.UserTokenRepositoryInterface
}
