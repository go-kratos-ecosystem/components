package scopes

import (
	"gorm.io/gorm"
)

// Offset add offset condition
//
//	Offset(3)
func Offset(offset int) *Scopes {
	return New().Offset(offset)
}

// Skip add offset condition
//
//	Skip(3)
func Skip(offset int) *Scopes {
	return New().Skip(offset)
}

// Limit add limit condition
//
//	Limit(3)
func Limit(limit int) *Scopes {
	return New().Limit(limit)
}

// Take add limit condition
//
//	Take(3)
func Take(limit int) *Scopes {
	return New().Take(limit)
}

// Page add page condition
//
//	Page(2, 10)
func Page(page, prePage int) *Scopes {
	return New().Page(page, prePage)
}

// Offset add offset condition
//
//	Offset(3)
func (s *Scopes) Offset(offset int) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	})
}

// Skip add offset condition
//
//	Skip(3)
func (s *Scopes) Skip(offset int) *Scopes {
	return s.Offset(offset)
}

// Limit add limit condition
//
//	Limit(3)
func (s *Scopes) Limit(limit int) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})
}

// Take add limit condition
//
//	Take(3)
func (s *Scopes) Take(limit int) *Scopes {
	return s.Limit(limit)
}

// Page add page condition
//
//	Page(2, 10)
func (s *Scopes) Page(page, prePage int) *Scopes {
	return s.Limit(prePage).Offset((page - 1) * prePage)
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
