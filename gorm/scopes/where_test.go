package scopes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhere_Where(t *testing.T) {
	users := []*User{
		GetUser("WhereUser1", GetUserOptions{}),
		GetUser("WhereUser2", GetUserOptions{}),
		GetUser("WhereUser3", GetUserOptions{}),
	}

	DB.Create(&users)

	var users1, users2, users3 []User
	DB.Scopes(Where("name in (?)", []string{"WhereUser1", "WhereUser2"}).Scope()).Find(&users1)
	assert.Len(t, users1, 2)

	DB.Scopes(Where("name in (?)", []string{"WhereUser1", "WhereUser4"}).Scope()).Find(&users2)
	assert.Len(t, users2, 1)

	DB.Scopes(Where("name = ?", "WhereUser3").Scope()).Find(&users3)
	assert.Len(t, users3, 1)
	assert.Equal(t, "WhereUser3", users3[0].Name)
}
