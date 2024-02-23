package scope

import "gorm.io/gorm"

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
