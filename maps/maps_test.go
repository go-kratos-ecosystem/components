package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaps(t *testing.T) {
	maps := Maps{}.Set("name", "Flc").When(true, func(maps Maps) Maps {
		return maps.Set("age", 18)
	}).When(false, func(maps Maps) Maps {
		return maps.Set("age", 20)
	}).Map(func(key string, value interface{}) (string, interface{}) {
		if key == "age" {
			return key, 21
		}

		return key, value
	})

	assert.Equal(t, "Flc", maps["name"])
	assert.Equal(t, 21, maps["age"])
	assert.True(t, maps.Has("name"))
	assert.Equal(t, 2, maps.Len())
	assert.True(t, func() bool {
		ok := true
		for _, key := range maps.Keys() {
			if key != "name" && key != "age" {
				ok = false
				break
			}
		}
		return ok
	}())
	assert.True(t, func() bool {
		ok := true
		for _, value := range maps.Values() {
			if value != "Flc" && value != 21 {
				ok = false
				break
			}
		}
		return ok
	}())

	assert.Equal(t, maps, maps.Clone())
	assert.NotSame(t, maps, maps.Clone())

	assert.Equal(t, map[string]interface{}{"name": "Flc", "age": 21}, maps.Maps())
	assert.Equal(t, map[string]interface{}{"name": "Flc", "age": 21}, maps.All())

	maps.Each(func(key string, value interface{}) {
		if key == "name" {
			maps["first_name"] = value
		}

		if key == "age" {
			maps["year"] = value
		}
	})

	assert.Equal(t, "Flc", maps["first_name"])

	maps.Delete("first_name")
	assert.Nil(t, maps["first_name"])

	maps.Merge(Maps{"year": "123"})
	assert.Equal(t, "123", maps["year"])

	maps.Unless(true, func(maps Maps) Maps {
		return maps.Set("sex", "woman")
	}).Unless(false, func(maps Maps) Maps {
		return maps.Set("sex", "man")
	})
	assert.Equal(t, "man", maps["sex"])
}

func TestMaps_Get(t *testing.T) {
	m := Maps{
		"name": "Flc",
	}
	name, ok := m.Get("name")
	assert.Equal(t, "Flc", name)
	assert.True(t, ok)

	assert.Equal(t, "Flc", m["name"])
	assert.Equal(t, "Flc", m.GetX("name"))

	// none
	name2, ok2 := m.Get("none")
	assert.Nil(t, name2)
	assert.False(t, ok2)

	assert.Equal(t, nil, m["none"])
	assert.Panics(t, func() {
		m.GetX("none")
	})
}
