package repository

import (
	"context"
	"errors"

	"github.com/ming-0x0/scaffold/shared/domainerror"
	"github.com/ming-0x0/scaffold/shared/transaction"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DefaultPage    int64 = 1
	DefaultPerPage int64 = 20
)

type Domain interface {
	Validate() error
}

type RepositoryInterface[D Domain] interface {
	FindByConditionsWithPagination(
		ctx context.Context,
		pageData map[string]int64,
		conditions map[string]any,
		scopes ...func(*gorm.DB) *gorm.DB,
	) ([]D, int64, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]any,
		scopes ...func(*gorm.DB) *gorm.DB,
	) (D, error)
	DeleteByConditions(
		ctx context.Context,
		conditions map[string]any,
	) error
	FindByConditions(
		ctx context.Context,
		conditions map[string]any,
		scopes ...func(*gorm.DB) *gorm.DB,
	) ([]D, error)
	Save(
		ctx context.Context,
		domain D,
	) error
	PreloadAssociations(ctx context.Context) func(*gorm.DB) *gorm.DB
}

type Repository[D Domain] struct {
	db     *gorm.DB
	logger *logrus.Logger
}

var _ RepositoryInterface[Domain] = (*Repository[Domain])(nil)

func New[D Domain](
	db *gorm.DB,
	logger *logrus.Logger,
) *Repository[D] {
	return &Repository[D]{
		db:     db,
		logger: logger,
	}
}

func (r *Repository[D]) DB(ctx context.Context) *gorm.DB {
	if tx, ok := transaction.TransactionFromContext(ctx); ok {
		return tx
	}

	return r.db.WithContext(ctx)
}

func (r *Repository[D]) pagination(pageData map[string]int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := DefaultPage
		if valPage, ok := pageData["page"]; ok && valPage > 0 {
			page = valPage
		}

		pageSize := DefaultPerPage
		if valPageSize, ok := pageData["limit"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

func (r *Repository[D]) FindByConditions(
	ctx context.Context,
	conditions map[string]any,
	scopes ...func(*gorm.DB) *gorm.DB,
) ([]D, error) {
	var domains []D
	if err := r.DB(ctx).Scopes(scopes...).Where(conditions).Find(&domains).Error; err != nil {
		return domains, domainerror.Wrap(domainerror.Internal, err)
	}

	for _, domain := range domains {
		err := domain.Validate()
		if err != nil {
			return []D{}, err
		}
	}

	return domains, nil
}

func (r *Repository[D]) TakeByConditions(
	ctx context.Context,
	conditions map[string]any,
	scopes ...func(*gorm.DB) *gorm.DB,
) (D, error) {
	var domain D
	err := r.DB(ctx).Scopes(scopes...).Where(conditions).Take(&domain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain, domainerror.Wrap(domainerror.NotFound, err)
		}

		return domain, domainerror.Wrap(domainerror.Internal, err)
	}

	err = domain.Validate()
	if err != nil {
		return domain, err
	}

	return domain, nil
}

func (r *Repository[D]) Save(
	ctx context.Context,
	domain D,
) error {
	err := domain.Validate()
	if err != nil {
		return err
	}

	err = r.DB(ctx).Save(&domain).Error
	if err != nil {
		return domainerror.Wrap(domainerror.Internal, err)
	}

	return nil
}

func (r *Repository[D]) DeleteByConditions(
	ctx context.Context,
	conditions map[string]any,
) error {
	var domain D
	err := r.DB(ctx).Where(conditions).Delete(&domain).Error
	if err != nil {
		return domainerror.Wrap(domainerror.Internal, err)
	}

	return nil
}

func (r *Repository[D]) FindByConditionsWithPagination(
	ctx context.Context,
	pageData map[string]int64,
	conditions map[string]any,
	scopes ...func(*gorm.DB) *gorm.DB,
) ([]D, int64, error) {
	cdb := r.DB(ctx)

	var domains []D
	var count int64

	countBuilder := cdb.Model(&domains)
	queryBuilder := cdb.Scopes(r.pagination(pageData))

	err := countBuilder.Scopes(scopes...).Where(conditions).Count(&count).Error
	if err != nil {
		return []D{}, 0, domainerror.Wrap(domainerror.Internal, err)
	}

	err = queryBuilder.Scopes(scopes...).Where(conditions).Find(&domains).Error
	if err != nil {
		return []D{}, 0, domainerror.Wrap(domainerror.Internal, err)
	}

	for _, domain := range domains {
		err := domain.Validate()
		if err != nil {
			return domains, count, err
		}
	}

	return domains, count, nil
}

func (r *Repository[D]) PreloadAssociations(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(clause.Associations)
	}
}
