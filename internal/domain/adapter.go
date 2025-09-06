package domain

type AdapterInterface interface {
	BannerRepository() BannerRepositoryInterface
	CategoryRepository() CategoryRepositoryInterface
	CustomerRepository() CustomerRepositoryInterface
	FooterRepository() FooterRepositoryInterface
	PartnerRepository() PartnerRepositoryInterface
	PermissionGroupRepository() PermissionGroupRepositoryInterface
	PermissionRepository() PermissionRepositoryInterface
	PostRepository() PostRepositoryInterface
	ResourceRepository() ResourceRepositoryInterface
	RoleRepository() RoleRepositoryInterface
	UserRepository() UserRepositoryInterface
	UserTokenRepository() UserTokenRepositoryInterface
}
