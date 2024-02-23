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

func WhereBetween(field string, start, end interface{}) *Scopes {
	return New().WhereBetween(field, start, end)
}

func WhereNotBetween(field string, start, end interface{}) *Scopes {
	return New().WhereNotBetween(field, start, end)
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
