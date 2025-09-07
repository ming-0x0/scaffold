package repository

import (
	"context"
	"errors"

	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"github.com/ming-0x0/scaffold/internal/shared/repository/condition"
	"github.com/ming-0x0/scaffold/internal/shared/transaction"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	DefaultPage    int32 = 1
	DefaultPerPage int32 = 20
)

type Domain interface {
	Validate() error
}

type RepositoryInterface[D Domain] interface {
	FindByConditionsWithPagination(
		ctx context.Context,
		pageData map[string]int32,
		conditions ...condition.Condition,
	) ([]D, int32, error)
	TakeByConditions(
		ctx context.Context,
		conditions ...condition.Condition,
	) (D, error)
	DeleteByConditions(
		ctx context.Context,
		conditions ...condition.Condition,
	) error
	FindByConditions(
		ctx context.Context,
		conditions ...condition.Condition,
	) ([]D, error)
	Save(
		ctx context.Context,
		domain D,
	) error
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

func (r *Repository[D]) scopes(conditions []condition.Condition) []func(*gorm.DB) *gorm.DB {
	scopes := make([]func(*gorm.DB) *gorm.DB, len(conditions))
	for i, condition := range conditions {
		scopes[i] = condition
	}

	return scopes
}

func (r *Repository[D]) pagination(pageData map[string]int32) func(*gorm.DB) *gorm.DB {
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
	conditions ...condition.Condition,
) ([]D, error) {
	var domains []D
	if err := r.DB(ctx).Scopes(r.scopes(conditions)...).Find(&domains).Error; err != nil {
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
	conditions ...condition.Condition,
) (D, error) {
	var domain D
	err := r.DB(ctx).Scopes(r.scopes(conditions)...).Take(&domain).Error
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
	conditions ...condition.Condition,
) error {
	var domain D
	err := r.DB(ctx).Scopes(r.scopes(conditions)...).Delete(&domain).Error
	if err != nil {
		return domainerror.Wrap(domainerror.Internal, err)
	}

	return nil
}

func (r *Repository[D]) FindByConditionsWithPagination(
	ctx context.Context,
	pageData map[string]int32,
	conditions ...condition.Condition,
) ([]D, int32, error) {
	cdb := r.DB(ctx)

	var domains []D
	var count int64

	countBuilder := cdb.Model(&domains)
	queryBuilder := cdb.Scopes(r.pagination(pageData))

	err := countBuilder.Scopes(r.scopes(conditions)...).Count(&count).Error
	if err != nil {
		return []D{}, 0, domainerror.Wrap(domainerror.Internal, err)
	}

	err = queryBuilder.Scopes(r.scopes(conditions)...).Find(&domains).Error
	if err != nil {
		return []D{}, 0, domainerror.Wrap(domainerror.Internal, err)
	}

	for _, domain := range domains {
		err := domain.Validate()
		if err != nil {
			return domains, 0, err
		}
	}

	return domains, int32(count), nil
}
