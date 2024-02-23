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

	CleanUsers()
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

func TestWhere_Between(t *testing.T) {
	users := []*User{
		GetUser("WhereBetweenUser1", GetUserOptions{Age: 18}),
		GetUser("WhereBetweenUser2", GetUserOptions{Age: 20}),
		GetUser("WhereBetweenUser3", GetUserOptions{Age: 22}),
	}

	CleanUsers()
	DB.Create(&users)

	var users1, users2, users3 []User
	DB.Scopes(WhereBetween("age", 18, 20).Scope()).Find(&users1)
	assert.Len(t, users1, 2)

	DB.Scopes(WhereBetween("age", 18, 19).Scope()).Find(&users2)
	assert.Len(t, users2, 1)
	assert.Equal(t, "WhereBetweenUser1", users2[0].Name)

	DB.Scopes(WhereBetween("age", 12, 16).Scope()).Find(&users3)
	assert.Len(t, users3, 0)

	var users4, users5, users6 []User
	DB.Scopes(WhereNotBetween("age", 18, 20).Scope()).Find(&users4)
	assert.Len(t, users4, 1)
	assert.Equal(t, "WhereBetweenUser3", users4[0].Name)

	DB.Scopes(WhereNotBetween("age", 18, 19).Scope()).Find(&users5)
	assert.Len(t, users5, 2)

	DB.Scopes(WhereNotBetween("age", 12, 16).Scope()).Find(&users6)
	assert.Len(t, users6, 3)
}

func TestWhere_In(t *testing.T) {
	users := []*User{
		GetUser("WhereInUser1", GetUserOptions{Age: 18}),
		GetUser("WhereInUser2", GetUserOptions{Age: 20}),
		GetUser("WhereInUser3", GetUserOptions{Age: 22}),
	}

	CleanUsers()
	DB.Create(&users)

	var users1, users2, users3 []User
	DB.Scopes(WhereIn("name", "WhereInUser1", "WhereInUser2").Scope()).Find(&users1)
	assert.Len(t, users1, 2)

	DB.Scopes(WhereIn("age", []int{18, 20}).Scope()).Find(&users2)
	assert.Len(t, users2, 2)

	DB.Scopes(WhereIn("name", []string{"WhereInUser1", "WhereInUser2"}).Scope()).Find(&users3)
	assert.Len(t, users3, 2)

	var users4, users5, users6 []User
	DB.Scopes(WhereNotIn("name", "WhereInUser1", "WhereInUser2").Scope()).Find(&users4)
	assert.Len(t, users4, 1)
	assert.Equal(t, "WhereInUser3", users4[0].Name)

	DB.Scopes(WhereNotIn("age", []int{18, 20}).Scope()).Find(&users5)
	assert.Len(t, users5, 1)
	assert.Equal(t, "WhereInUser3", users5[0].Name)

	DB.Scopes(WhereNotIn("name", []string{"WhereInUser1", "WhereInUser2"}).Scope()).Find(&users6)
	assert.Len(t, users6, 1)
	assert.Equal(t, "WhereInUser3", users6[0].Name)
}

func TestWhere_Like(t *testing.T) {
	users := []*User{
		GetUser("WhereLikeUser1", GetUserOptions{Age: 18}),
		GetUser("WhereLikeUser2", GetUserOptions{Age: 20}),
		GetUser("WhereLikeUser3", GetUserOptions{Age: 22}),
	}

	CleanUsers()
	DB.Create(&users)

	var users1, users2, users3, users4 []User
	DB.Scopes(WhereLike("name", "WhereLikeUser1").Scope()).Find(&users1)
	assert.Len(t, users1, 1)
	assert.Equal(t, "WhereLikeUser1", users1[0].Name)

	DB.Scopes(WhereLike("name", "WhereLike%").Scope()).Find(&users2)
	assert.Len(t, users2, 3)

	DB.Scopes(WhereLike("name", "%LikeUser3").Scope()).Find(&users3)
	assert.Len(t, users3, 1)
	assert.Equal(t, "WhereLikeUser3", users3[0].Name)

	DB.Scopes(WhereLike("name", "%Like%").Scope()).Find(&users4)
	assert.Len(t, users4, 3)

	var users5, users6, users7, users8 []User
	DB.Scopes(WhereNotLike("name", "WhereLikeUser1").Scope()).Find(&users5)
	assert.Len(t, users5, 2)

	DB.Scopes(WhereNotLike("name", "WhereLike%").Scope()).Find(&users6)
	assert.Len(t, users6, 0)

	DB.Scopes(WhereNotLike("name", "%LikeUser3").Scope()).Find(&users7)
	assert.Len(t, users7, 2)

	DB.Scopes(WhereNotLike("name", "%Like%").Scope()).Find(&users8)
	assert.Len(t, users8, 0)
}
