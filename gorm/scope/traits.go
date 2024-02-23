package scope

import "gorm.io/gorm"

func (s *Scopes) When(condition bool, fc func(*gorm.DB) *gorm.DB) *Scopes {
	if condition {
		return s.Add(fc)
	}
	return s
}

func (s *Scopes) Unless(condition bool, fc func(*gorm.DB) *gorm.DB) *Scopes {
	return s.When(!condition, fc)
}
