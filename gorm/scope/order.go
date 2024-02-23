package scope

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func (s *Scopes) OrderBy(column string, reorder ...string) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", column, s.buildReorder(reorder...)))
	})
}

func (s *Scopes) OrderByDesc(column string) *Scopes {
	return s.OrderBy(column, "desc")
}

func (s *Scopes) OrderByAsc(column string) *Scopes {
	return s.OrderBy(column, "asc")
}

func (s *Scopes) OrderByRaw(sql string, values ...interface{}) *Scopes {
	return s.Add(func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf(sql, values...)) // TODO: 语法待定？
	})
}

func (s *Scopes) buildReorder(reorder ...string) string {
	if len(reorder) > 0 && strings.ToUpper(reorder[0]) == "DESC" {
		return "DESC"
	}
	return "ASC"
}
