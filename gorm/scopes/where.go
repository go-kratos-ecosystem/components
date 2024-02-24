package scopes

import (
	"fmt"

	"gorm.io/gorm"
)

// Where add where condition
//
//	Where("name = ?", "Flc")
//	Where("name = ? AND age = ?", "Flc", 20)
func Where(query interface{}, args ...interface{}) *Scopes {
	return New().Where(query, args...)
}

// WhereBetween add where between condition
//
//	WhereBetween("age", 18, 20)
func WhereBetween(field string, start, end interface{}) *Scopes {
	return New().WhereBetween(field, start, end)
}

// WhereNotBetween add where not between condition
//
//	WhereNotBetween("age", 18, 20)
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

// WhereLike add where like condition
//
//	WhereLike("name", "Flc")
//	WhereLike("name", "Flc%")
//	WhereLike("name", "%Flc")
//	WhereLike("name", "%Flc%")
func WhereLike(field string, value interface{}) *Scopes {
	return New().WhereLike(field, value)
}

// WhereNotLike add where not like condition
//
//	WhereNotLike("name", "Flc")
//	WhereNotLike("name", "Flc%")
//	WhereNotLike("name", "%Flc")
//	WhereNotLike("name", "%Flc%")
func WhereNotLike(field string, value interface{}) *Scopes {
	return New().WhereNotLike(field, value)
}

// WhereEq add where eq condition
//
//	WhereEq("name", "Flc")
//	WhereEq("age", 18)
func WhereEq(field string, value interface{}) *Scopes {
	return New().WhereEq(field, value)
}

// WhereEgt add where egt condition
//
//	WhereEgt("age", 18)
func WhereEgt(field string, value interface{}) *Scopes {
	return New().WhereEgt(field, value)
}

// WhereGt add where gt condition
//
//	WhereGt("age", 18)
func WhereGt(field string, value interface{}) *Scopes {
	return New().WhereGt(field, value)
}

// WhereElt add where elt condition
//
//	WhereElt("age", 18)
func WhereElt(field string, value interface{}) *Scopes {
	return New().WhereElt(field, value)
}

// WhereLt add where lt condition
//
//	WhereLt("age", 18)
func WhereLt(field string, value interface{}) *Scopes {
	return New().WhereLt(field, value)
}

// WhereNe add where ne condition
//
//	WhereNe("name", "Flc")
//	WhereNe("age", 18)
func WhereNe(field string, value interface{}) *Scopes {
	return New().WhereNe(field, value)
}

// WhereNot add where not condition
//
//	WhereNot("name = ?", "Flc")
//	WhereNot("name = ? AND age = ?", "Flc", 20)
func WhereNot(query interface{}, args ...interface{}) *Scopes {
	return New().WhereNot(query, args...)
}

// WhereNull add where null condition
//
//	WhereNull("name")
func WhereNull(field string) *Scopes {
	return New().WhereNull(field)
}

// WhereNotNull add where not null condition
//
//	WhereNotNull("name")
func WhereNotNull(field string) *Scopes {
	return New().WhereNotNull(field)
}

// Where add where condition
//
//	Where("name = ?", "Flc")
//	Where("name = ? AND age = ?", "Flc", 20)
func (s *Scopes) Where(query interface{}, args ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	})
}

// WhereBetween add where between condition
//
//	WhereBetween("age", 18, 20)
func (s *Scopes) WhereBetween(column string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), start, end)
	})
}

// WhereNotBetween add where not between condition
//
//	WhereNotBetween("age", 18, 20)
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

// WhereLike add where like condition
//
//	WhereLike("name", "Flc")
//	WhereLike("name", "Flc%")
//	WhereLike("name", "%Flc")
//	WhereLike("name", "%Flc%")
func (s *Scopes) WhereLike(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s LIKE ?", column), value)
	})
}

// WhereNotLike add where not like condition
//
//	WhereNotLike("name", "Flc")
//	WhereNotLike("name", "Flc%")
//	WhereNotLike("name", "%Flc")
//	WhereNotLike("name", "%Flc%")
func (s *Scopes) WhereNotLike(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s NOT LIKE ?", column), value)
	})
}

// WhereEq add where eq condition
//
//	WhereEq("name", "Flc")
//	WhereEq("age", 18)
func (s *Scopes) WhereEq(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", column), value)
	})
}

// WhereEgt add where egt condition
//
//	WhereEgt("age", 18)
func (s *Scopes) WhereEgt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s >= ?", column), value)
	})
}

// WhereGt add where gt condition
//
//	WhereGt("age", 18)
func (s *Scopes) WhereGt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s > ?", column), value)
	})
}

// WhereElt add where elt condition
//
//	WhereElt("age", 18)
func (s *Scopes) WhereElt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <= ?", column), value)
	})
}

// WhereLt add where lt condition
//
//	WhereLt("age", 18)
func (s *Scopes) WhereLt(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s < ?", column), value)
	})
}

// WhereNe add where ne condition
//
//	WhereNe("name", "Flc")
//	WhereNe("age", 18)
func (s *Scopes) WhereNe(column string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <> ?", column), value)
	})
}

// WhereNot add where not condition
//
//	WhereNot("name = ?", "Flc")
//	WhereNot("name = ? AND age = ?", "Flc", 20)
func (s *Scopes) WhereNot(query interface{}, args ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Not(query, args...)
	})
}

// WhereNull add where null condition
//
//	WhereNull("name")
func (s *Scopes) WhereNull(column string) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NULL", column))
	})
}

// WhereNotNull add where not null condition
//
//	WhereNotNull("name")
func (s *Scopes) WhereNotNull(column string) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NOT NULL", column))
	})
}
