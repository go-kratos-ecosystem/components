package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Name string
	Age  int
}

func TestTap_Struct(t *testing.T) {
	f := &foo{Name: "foo"}

	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	f = Tap(f, func(f *foo) {
		f.Name = "bar" //nolint:goconst
		f.Age = 18
	})
	assert.Equal(t, "bar", f.Name)
	assert.Equal(t, 18, f.Age)
}

func TestTap_Int(t *testing.T) {
	f := new(int)
	*f = 10

	assert.Equal(t, 10, *f)
	f = Tap(f, func(f *int) {
		*f = 20
	})
	assert.Equal(t, 20, *f)

	b := 10
	assert.Equal(t, 10, b)
	b = Tap(b, func(b int) { //nolint:staticcheck
		b = 20 //nolint:ineffassign,staticcheck
		_ = b
	})
	assert.Equal(t, 10, b)

	b2 := Tap(&b, func(b *int) {
		*b = 20
	})
	assert.Equal(t, 20, *b2)
}

func BenchmarkTap_UseTap(b *testing.B) {
	f := &foo{Name: "foo"}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Tap(f, func(f *foo) {
				f.Name = "bar"
				f.Age = 18
			})
		}
	})
}

func BenchmarkTap_UnUseTap(b *testing.B) {
	f := &foo{Name: "foo"}
	fc := func(f *foo, callbacks ...func(*foo)) {
		for _, callback := range callbacks {
			if callback != nil {
				callback(f)
			}
		}
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fc(f, func(f *foo) {
				f.Name = "bar"
				f.Age = 18
			})
		}
	})
}

func TestWith(t *testing.T) {
	f := &foo{Name: "foo"}

	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	f2 := With(f, func(f *foo) *foo {
		f.Name = "bar"
		f.Age = 18
		return f
	})
	assert.Equal(t, "bar", f2.Name)
	assert.Equal(t, 18, f2.Age)
}

func TestWhen(t *testing.T) {
	f := &foo{Name: "foo"}

	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	f2 := When(f, true, func(f *foo) *foo {
		f.Name = "bar"
		f.Age = 18
		return f
	})
	assert.Equal(t, "bar", f2.Name)
	assert.Equal(t, 18, f2.Age)

	f3 := When(f, false, func(f *foo) *foo {
		f.Name = "baz" //nolint:goconst
		f.Age = 20
		return f
	})
	assert.Equal(t, "bar", f3.Name)
	assert.Equal(t, 18, f3.Age)
}
