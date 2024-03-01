package maps

import (
	"fmt"
)

type M map[string]interface{}

func (m M) Maps() map[string]interface{} {
	return m
}

func (m M) All() map[string]interface{} {
	return m.Maps()
}

func (m M) Merge(n M) M {
	for k, v := range n {
		m[k] = v
	}
	return m
}

func (m M) Clone() M {
	n := make(M, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}

func (m M) Has(k string) bool {
	_, ok := m[k]
	return ok
}

func (m M) Get(k string) (interface{}, bool) {
	v, ok := m[k]
	return v, ok
}

func (m M) GetX(k string) interface{} {
	if v, ok := m.Get(k); ok {
		return v
	}

	panic(fmt.Sprintf("maps: key %s not exists", k))
}

func (m M) Set(k string, v interface{}) M {
	m[k] = v
	return m
}

func (m M) Delete(k string) M {
	delete(m, k)
	return m
}

func (m M) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m M) Values() []interface{} {
	values := make([]interface{}, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func (m M) Len() int {
	return len(m)
}

func (m M) When(guard bool, fn func(maps M) M) M {
	if guard {
		return fn(m)
	}
	return m
}

func (m M) Unless(guard bool, fn func(maps M) M) M {
	if !guard {
		return fn(m)
	}
	return m
}

func (m M) Map(fn func(k string, v interface{}) (string, interface{})) M {
	n := make(M, len(m))
	for k, v := range m {
		k, v := fn(k, v)
		n[k] = v
	}
	return n
}

func (m M) Each(fn func(k string, v interface{})) {
	for k, v := range m {
		fn(k, v)
	}
}
