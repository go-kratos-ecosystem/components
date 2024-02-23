package scopes

import "gorm.io/gorm"

func Where(query interface{}, args ...interface{}) *Scopes {
	return New().Where(query, args...)
}

func WhereBetween(field string, start, end interface{}) *Scopes {
	return New().WhereBetween(field, start, end)
}

func WhereNotBetween(field string, start, end interface{}) *Scopes {
	return New().WhereNotBetween(field, start, end)
}

func (s *Scopes) Where(query interface{}, args ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	})
}

func (s *Scopes) WhereBetween(column string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? BETWEEN ? AND ?", column, start, end)
	})
}

func (s *Scopes) WhereNotBetween(column string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? NOT BETWEEN ? AND ?", column, start, end)
	})
}

func (s *Scopes) WhereIn(column string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? IN (?)", column, values)
	})
}

func (s *Scopes) WhereNotIn(column string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? NOT IN (?)", column, values)
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
