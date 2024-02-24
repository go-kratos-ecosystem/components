package scopes

import (
	"gorm.io/gorm"
)

func (s *Scopes) Skip(offset int) *Scopes {
	return s.Offset(offset)
}

func (s *Scopes) Offset(offset int) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	})
}

func (s *Scopes) Limit(limit int) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})
}

func (s *Scopes) Take(limit int) *Scopes {
	return s.Limit(limit)
}

func (s *Scopes) Page(page, prePage int) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * prePage).Limit(prePage)
	})
}

// Paginate TODO: 未实现
// func (s *Scopes) Paginate(page, pageSize int, dest interface{}) *Scopes {
// 	var total int64
//
// 	dest = pagination.Paginator{
// 		Page:    page,
// 		PrePage: pageSize,
// 		Total:   int(total),
// 	}
//
// 	return s.Add(func(db *gorm.DB) *gorm.DB {
// 		db.Count(&total)
// 		return db.Offset((page - 1) * pageSize).Limit(pageSize)
// 	})
// }
