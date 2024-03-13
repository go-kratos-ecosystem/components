package helper

import (
	"context"
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

func TestChain(t *testing.T) {
	// chain functions
	chain := Chain(
		func(s string) string {
			return s + "1"
		},
		func(s string) string {
			return s + "2"
		},
		func(s string) string {
			return s + "3"
		},
	)

	assert.Equal(t, "0123", chain("0"))

	// chain functions
	chain2 := Chain(
		func(foo *foo) *foo {
			foo.Name = "bar"
			return foo
		},
		func(foo *foo) *foo {
			foo.Age = 18
			return foo
		},
	)

	f := &foo{Name: "foo"}
	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	got := chain2(f)
	assert.Equal(t, "bar", got.Name)
	assert.Equal(t, 18, got.Age)
}

func TestScan_Basic(t *testing.T) {
	// string
	var foo string
	err := Scan("foo", &foo)
	assert.Nil(t, err)
	assert.Equal(t, "foo", foo)

	// int
	var bar int
	err = Scan(1, &bar)
	assert.Nil(t, err)
	assert.Equal(t, 1, bar)

	// struct
	type Baz struct {
		Name string
	}

	// struct.1
	var baz Baz
	err = Scan(Baz{
		Name: "baz",
	}, &baz)
	assert.Nil(t, err)
	assert.Equal(t, "baz", baz.Name)

	// struct.2
	assert.Error(t, Scan("foo", nil))
	var baz2 *Baz
	assert.Error(t, Scan("foo", baz2))

	// struct.3
	var baz3 Baz
	err = Scan(func() interface{} {
		return Baz{
			Name: "baz",
		}
	}(), &baz3)
	assert.Nil(t, err)
	assert.Equal(t, "baz", baz3.Name)

	// test lower
	type test struct {
		Name string
	}
	var tt test
	err = Scan(func() interface{} {
		return test{
			Name: "test",
		}
	}(), &tt)
	assert.Nil(t, err)
	assert.Equal(t, "test", tt.Name)
}

func TestScan_ComplexStruct(t *testing.T) {
	type AName struct {
		Name string
	}

	type ACompany struct {
		Name string
	}

	type A struct {
		Name      *AName
		Companies []*ACompany
	}

	type BName struct {
		Name string
	}

	type BCompany struct {
		Name string
	}

	type B struct {
		Name      *BName
		Companies []*BCompany
	}

	a := &A{
		Name: &AName{
			Name: "A",
		},
		Companies: []*ACompany{
			{
				Name: "A1",
			},
			{
				Name: "A2",
			},
		},
	}

	var b B
	err := Scan(a, &b)
	assert.Nil(t, err)
	assert.Equal(t, "A", b.Name.Name)
	assert.Equal(t, "A1", b.Companies[0].Name)
	assert.Equal(t, "A2", b.Companies[1].Name)
}

func TestChainWithErr(t *testing.T) {
	// chain functions
	chain := ChainWithErr(
		func(s string) (string, error) {
			return s + "1", nil
		},
		func(s string) (string, error) {
			return s + "2", nil
		},
		func(s string) (string, error) {
			return s + "3", nil
		},
	)

	got, err := chain("0")
	assert.Nil(t, err)
	assert.Equal(t, "0123", got)

	// chain functions
	chain2 := ChainWithErr(
		func(foo *foo) (*foo, error) {
			foo.Name = "bar"
			return foo, nil
		},
		func(foo *foo) (*foo, error) {
			foo.Age = 18
			return foo, nil
		},
	)

	f := &foo{Name: "foo"}
	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	got2, err := chain2(f)
	assert.Nil(t, err)
	assert.Equal(t, "bar", got2.Name)
	assert.Equal(t, 18, got2.Age)

	// context
	chain3 := ChainWithErr(
		func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, "foo", "bar"), nil //nolint:revive,staticcheck
		},
		func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, "bar", "baz"), nil //nolint:revive,staticcheck
		},
	)

	ctx, err := chain3(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "bar", ctx.Value("foo"))
	assert.Equal(t, "baz", ctx.Value("bar"))

	// context with error
	chain4 := ChainWithErr(
		func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, "foo", "bar"), nil //nolint:revive,staticcheck
		},
		func(context.Context) (context.Context, error) {
			return nil, assert.AnError
		},
	)

	ctx, err = chain4(context.Background())
	assert.Error(t, err)
	assert.Nil(t, ctx)
}

func TestIf(t *testing.T) {
	// if true
	got := If(true, "foo", "bar")
	assert.Equal(t, "foo", got)

	// if false
	got = If(false, "foo", "bar")
	assert.Equal(t, "bar", got)
}
