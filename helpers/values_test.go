package helpers

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestTransform(t *testing.T) {
	got := Transform("foo", func(s string) *foo {
		return &foo{Name: s}
	})
	assert.Equal(t, "foo", got.Name)

	assert.Equal(t, "1", Transform(1, strconv.Itoa))
}

func TestIf(t *testing.T) {
	// if true
	got := If(true, "foo", "bar")
	assert.Equal(t, "foo", got)

	// if false
	got = If(false, "foo", "bar")
	assert.Equal(t, "bar", got)

	// slices
	opts := []string{"foo", "bar"}
	got2 := If(true, append(opts, "baz"), opts)
	assert.Equal(t, []string{"foo", "bar", "baz"}, got2)
	got3 := If(false, append(opts, "baz"), opts)
	assert.Equal(t, []string{"foo", "bar"}, got3)

	// pointers
	f := &foo{Name: "foo"}
	got4 := If(true, &foo{Name: "bar"}, f)
	assert.Equal(t, "bar", got4.Name)
	got5 := If(false, &foo{Name: "bar"}, f)
	assert.Equal(t, "foo", got5.Name)

	// pointers is nil
	var nilVal *foo
	assert.Panics(t, func() {
		If(nilVal != nil, nilVal.Name, "")
	})
}

func TestIfFunc(t *testing.T) {
	// if true
	got := IfFunc(true, func() string {
		return "foo" //nolint:goconst
	}, func() string {
		return "bar"
	})
	assert.Equal(t, "foo", got)

	// if false
	got = IfFunc(false, func() string {
		return "foo"
	}, func() string {
		return "bar"
	})
	assert.Equal(t, "bar", got)

	// zero pointer
	var nilVal *foo
	got = IfFunc(nilVal != nil, func() string {
		return nilVal.Name
	}, func() string {
		return "nil foo"
	})
	assert.Equal(t, "nil foo", got)
}

func TestUnless(t *testing.T) {
	// unless true
	got := Unless(true, "foo", "bar")
	assert.Equal(t, "bar", got)

	// unless false
	got = Unless(false, "foo", "bar")
	assert.Equal(t, "foo", got)
}

func TestOptional(t *testing.T) {
	// not nil
	got1 := Optional(&foo{Name: "bar"})
	assert.Equal(t, "bar", got1.Name)

	// nil
	got2 := Optional[foo](nil)
	assert.Equal(t, "", got2.Name)

	// nil val
	var nilVal *foo
	got3 := Optional(nilVal)
	assert.Equal(t, "", got3.Name)

	// int ptr
	got4 := Optional[int](nil)
	assert.Equal(t, 0, *got4)
}

func TestDefault(t *testing.T) {
	// string
	got := Default("", "foo")
	assert.Equal(t, "foo", got)

	// int
	got2 := Default(0, 10)
	assert.Equal(t, 10, got2)

	// struct
	got3 := Default(foo{}, foo{Name: "bar"})
	assert.Equal(t, "bar", got3.Name)

	// ptr
	got4 := Default(nil, &foo{Name: "bar"})
	assert.Equal(t, "bar", got4.Name)
}

func TestDefaultWith(t *testing.T) {
	// string
	got := DefaultWithFunc("", func() string {
		return "foo"
	})
	assert.Equal(t, "foo", got)

	// int
	got2 := DefaultWithFunc(0, func() int {
		return 10
	})
	assert.Equal(t, 10, got2)

	// struct
	got3 := DefaultWithFunc(foo{}, func() foo {
		return foo{Name: "bar"}
	})
	assert.Equal(t, "bar", got3.Name)

	// ptr
	got4 := DefaultWithFunc(nil, func() *foo {
		return &foo{Name: "bar"}
	})
	assert.Equal(t, "bar", got4.Name)
}

func TestPtrAndVal(t *testing.T) {
	// string
	got := Ptr("foo")
	assert.Equal(t, "foo", *got)
	assert.Equal(t, "foo", Val(got))

	// int
	got2 := Ptr(10)
	assert.Equal(t, 10, *got2)
	assert.Equal(t, 10, Val(got2))

	// struct
	got3 := Ptr(foo{Name: "bar"})
	assert.Equal(t, "bar", got3.Name)
	assert.Equal(t, "bar", Val(got3).Name)

	// time.Time
	now := time.Now()
	got4 := Ptr(now)
	assert.Equal(t, now.String(), got4.String())
	assert.Equal(t, now.String(), Val(got4).String())

	// nil
	got5 := Ptr[*int](nil)
	assert.Nil(t, *got5)

	// nil val
	var nilVal *int
	got6 := Val(nilVal)
	assert.Equal(t, 0, got6)
}

type testInterface interface {
	testFunc()
}

type testStruct struct{}

func (t testStruct) testFunc() {}

func TestIsType(t *testing.T) {
	assert.True(t, IsType[int](10))
	assert.False(t, IsType[int](int8(10)))

	assert.True(t, IsType[string]("foo"))
	assert.False(t, IsType[string](10))

	assert.True(t, IsType[time.Time](time.Now()))
	assert.False(t, IsType[time.Time](10))

	assert.True(t, IsType[foo](foo{}))
	assert.False(t, IsType[foo](10))

	assert.True(t, IsType[*foo](&foo{}))
	assert.False(t, IsType[foo](&foo{}))
	assert.False(t, IsType[*foo](nil))

	assert.True(t, IsType[testInterface](testStruct{}))
	assert.True(t, IsType[interface{}](testStruct{}))
	assert.True(t, IsType[any](testStruct{}))

	assert.True(t, IsType[error](errors.New("foo")))
	assert.False(t, IsType[error](nil))
}

func testIsType(value any) bool {
	_, ok := value.(int)
	return ok
}

func BenchmarkIsType_Native(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testIsType(10)
		}
	})
}

func BenchmarkIsType_Generics(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IsType[int](10)
		}
	})
}

func TestIsZero(t *testing.T) {
	assert.True(t, IsZero(0))
	assert.True(t, IsZero(""))
	assert.True(t, IsZero(false))
	assert.True(t, IsZero[*foo](nil))
	assert.True(t, IsZero(0.0))
	assert.True(t, IsZero(0.0+0i))

	assert.False(t, IsZero(1))
	assert.False(t, IsZero("foo"))
	assert.False(t, IsZero(true))
	assert.False(t, IsZero(1.0))
	assert.False(t, IsZero(1.0+0i))
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, IsEmpty(""))
	assert.True(t, IsEmpty[*foo](nil))
	assert.True(t, IsEmpty(0))
	assert.True(t, IsEmpty(false))
	assert.True(t, IsEmpty(0.0))
	assert.True(t, IsEmpty(0.0+0i))

	assert.False(t, IsEmpty("foo"))
	assert.False(t, IsEmpty(1))
	assert.False(t, IsEmpty(true))
	assert.False(t, IsEmpty(1.0))
	assert.False(t, IsEmpty(1.0+0i))
}
