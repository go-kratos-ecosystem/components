package scopes

import "gorm.io/gorm"

func (s *Scopes) Select(query interface{}, args ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Select(query, args...)
	})
}
