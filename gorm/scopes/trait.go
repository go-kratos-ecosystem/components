package scopes

import "gorm.io/gorm"

func When(condition bool, f func(db *gorm.DB) *gorm.DB) *Scopes {
	return New().When(condition, f)
}

func Unless(condition bool, f func(db *gorm.DB) *gorm.DB) *Scopes {
	return New().Unless(condition, f)
}

func (s *Scopes) When(condition bool, fc func(*gorm.DB) *gorm.DB) *Scopes {
	if condition {
		return s.Add(fc)
	}
	return s
}

func (s *Scopes) Unless(condition bool, fc func(*gorm.DB) *gorm.DB) *Scopes {
	return s.When(!condition, fc)
}
