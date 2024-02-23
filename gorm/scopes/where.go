package scopes

import (
	"fmt"

	"gorm.io/gorm"
)

func Where(query interface{}, args ...interface{}) *Scopes {
	return New().Where(query, args...)
}

func WhereBetween(field string, start, end interface{}) *Scopes {
	return New().WhereBetween(field, start, end)
}

func WhereNotBetween(field string, start, end interface{}) *Scopes {
	return New().WhereNotBetween(field, start, end)
}

// WhereIn add where in condition
//
//	WhereIn("name", []string{"WhereInUser1", "WhereInUser2"})
//	WhereIn("age", []int{18, 20})
//	WhereIn("name", "WhereInUser1", "WhereInUser2")
func WhereIn(field string, values ...interface{}) *Scopes {
	return New().WhereIn(field, values...)
}

// WhereNotIn add where not in condition
//
//	WhereNotIn("name", []string{"WhereInUser1", "WhereInUser2"})
//	WhereNotIn("age", []int{18, 20})
//	WhereNotIn("name", "WhereInUser1", "WhereInUser2")
func WhereNotIn(field string, values ...interface{}) *Scopes {
	return New().WhereNotIn(field, values...)
}

func WhereLike(field string, value interface{}) *Scopes {
	return New().WhereLike(field, value)
}

func WhereNotLike(field string, value interface{}) *Scopes {
	return New().WhereNotLike(field, value)
}

func WhereEq(field string, value interface{}) *Scopes {
	return New().WhereEq(field, value)
}

func WhereEgt(field string, value interface{}) *Scopes {
	return New().WhereEgt(field, value)
}

func WhereGt(field string, value interface{}) *Scopes {
	return New().WhereGt(field, value)
}

func WhereElt(field string, value interface{}) *Scopes {
	return New().WhereElt(field, value)
}

func WhereLt(field string, value interface{}) *Scopes {
	return New().WhereLt(field, value)
}

func (s *Scopes) Where(query interface{}, args ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	})
}

func (s *Scopes) WhereBetween(column string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), start, end)
	})
}

func (s *Scopes) WhereNotBetween(column string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), start, end)
	})
}

// WhereIn add where in condition
//
//	WhereIn("name", []string{"WhereInUser1", "WhereInUser2"})
//	WhereIn("age", []int{18, 20})
//	WhereIn("name", "WhereInUser1", "WhereInUser2")
func (s *Scopes) WhereIn(column string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		if len(values) > 1 {
			return db.Where(fmt.Sprintf("%s IN (?)", column), values)
		}
		return db.Where(fmt.Sprintf("%s IN ?", column), values...)
	})
}

// WhereNotIn add where not in condition
//
//	WhereNotIn("name", []string{"WhereInUser1", "WhereInUser2"})
//	WhereNotIn("age", []int{18, 20})
//	WhereNotIn("name", "WhereInUser1", "WhereInUser2")
func (s *Scopes) WhereNotIn(column string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		if len(values) > 1 {
			return db.Where(fmt.Sprintf("%s NOT IN (?)", column), values)
		}
		return db.Where(fmt.Sprintf("%s NOT IN ?", column), values...)
	})
}

func (s *Scopes) WhereLike(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? LIKE ?", column, value)
	})
}

func (s *Scopes) WhereNotLike(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? NOT LIKE ?", column, value)
	})
}

func (s *Scopes) WhereEq(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? = ?", column, value)
	})
}

func (s *Scopes) WhereEgt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? >= ?", column, value)
	})
}

func (s *Scopes) WhereGt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? > ?", column, value)
	})
}

func (s *Scopes) WhereElt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? <= ?", column, value)
	})
}

func (s *Scopes) WhereLt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? < ?", column, value)
	})
}

func (s *Scopes) WhereNe(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? <> ?", column, value)
	})
}

func (s *Scopes) WhereNull(column string) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? IS NULL", column)
	})
}

func (s *Scopes) WhereNotNull(column string) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? IS NOT NULL", column)
	})
}
