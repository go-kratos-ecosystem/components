package scope

import "gorm.io/gorm"

type Scopes []func(*gorm.DB) *gorm.DB

func New() *Scopes {
	return &Scopes{}
}

func When(condition bool, f func(db *gorm.DB) *gorm.DB) *Scopes {
	return New().When(condition, f)
}

func Unless(condition bool, f func(db *gorm.DB) *gorm.DB) *Scopes {
	return New().Unless(condition, f)
}

func Between(field string, start, end interface{}) *Scopes {
	return New().Between(field, start, end)
}

func NotBetween(field string, start, end interface{}) *Scopes {
	return New().NotBetween(field, start, end)
}

func (s *Scopes) Apply(db *gorm.DB) *gorm.DB {
	return db.Scopes(*s...)
}

func (s *Scopes) Scope() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return s.Apply(db)
	}
}

func (s *Scopes) Scopes() []func(*gorm.DB) *gorm.DB {
	return *s
}

func (s *Scopes) Add(scopes ...func(*gorm.DB) *gorm.DB) *Scopes {
	*s = append(*s, scopes...)
	return s
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

func (s *Scopes) Between(field string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? BETWEEN ? AND ?", field, start, end)
	})
}

func (s *Scopes) NotBetween(field string, start, end interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? NOT BETWEEN ? AND ?", field, start, end)
	})
}

func (s *Scopes) In(field string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? IN (?)", field, values)
	})
}

func (s *Scopes) NotIn(field string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? NOT IN (?)", field, values)
	})
}

func (s *Scopes) Like(field string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? LIKE ?", field, value)
	})
}

func (s *Scopes) NotLike(field string, value interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Where("? NOT LIKE ?", field, value)
	})
}
