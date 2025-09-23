package condition

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Condition func(*gorm.DB) *gorm.DB

func PreloadAssociations() Condition {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(clause.Associations)
	}
}

func EQ(column string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", column), value)
	}
}

func NEQ(column string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <> ?", column), value)
	}
}

func LIKE(column string, value any) Condition {
	return func(db *gorm.DB) *gorm.DB {
		trimmed := strings.TrimSpace(fmt.Sprint(value))
		escapedValue := "%" + strings.NewReplacer(
			"\\", "\\\\",
			"%", "\\%",
			"_", "\\_",
		).Replace(trimmed) + "%"

		return db.Where(
			fmt.Sprintf("LOWER(%s) LIKE LOWER(?) ESCAPE '\\\\'", column),
			escapedValue,
		)
	}
}

func OR(conditions ...Condition) Condition {
	return func(db *gorm.DB) *gorm.DB {
		if len(conditions) == 0 {
			return db
		}

		query := conditions[0](db)

		for _, cond := range conditions[1:] {
			query = query.Or(cond(db))
		}

		return query
	}
}

func IN(column string, values []any) Condition {
	if len(values) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN (?)", column), values)
	}
}

func IsNull(column string) Condition {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NULL", column))
	}
}
