package maps

import (
	"fmt"
)

type Maps map[string]interface{}

func (m Maps) Maps() map[string]interface{} {
	return m
}

func (m Maps) All() map[string]interface{} {
	return m.Maps()
}

func (m Maps) Merge(n Maps) Maps {
	for k, v := range n {
		m[k] = v
	}
	return m
}

func (m Maps) Clone() Maps {
	n := make(Maps, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}

func (m Maps) Has(k string) bool {
	_, ok := m[k]
	return ok
}

func (m Maps) Get(k string) (interface{}, bool) {
	v, ok := m[k]
	return v, ok
}

func (m Maps) GetX(k string) interface{} {
	if v, ok := m.Get(k); ok {
		return v
	}

	panic(fmt.Sprintf("maps: key %s not exists", k))
}

func (m Maps) Set(k string, v interface{}) Maps {
	m[k] = v
	return m
}

func (m Maps) Delete(k string) Maps {
	delete(m, k)
	return m
}

func (m Maps) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m Maps) Values() []interface{} {
	values := make([]interface{}, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func (m Maps) Len() int {
	return len(m)
}

func (m Maps) When(guard bool, fn func(maps Maps) Maps) Maps {
	if guard {
		return fn(m)
	}
	return m
}

func (m Maps) Unless(guard bool, fn func(maps Maps) Maps) Maps {
	if !guard {
		return fn(m)
	}
	return m
}

func (m Maps) Map(fn func(k string, v interface{}) (string, interface{})) Maps {
	n := make(Maps, len(m))
	for k, v := range m {
		k, v := fn(k, v)
		n[k] = v
	}
	return n
}

func (m Maps) Each(fn func(k string, v interface{})) {
	for k, v := range m {
		fn(k, v)
	}
}
