package scopes

import (
	"gorm.io/gorm"

	"github.com/go-kratos-ecosystem/components/v2/pagination"
)

func (s *Scopes) Page(page, pageSize int) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	})
}

// Paginate TODO: 未实现
func (s *Scopes) Paginate(page, pageSize int, dest interface{}) *Scopes {
	var total int64

	dest = pagination.Paginator{
		Page:    page,
		PrePage: pageSize,
		Total:   int(total),
	}

	return s.Add(func(db *gorm.DB) *gorm.DB {
		db.Count(&total)
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	})
}
