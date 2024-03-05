package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type proxyMock struct {
	Name string
	Age  int
}

func TestProxy_Point(t *testing.T) {
	mock := &proxyMock{Name: "foo"}

	proxy := NewProxy(mock)
	assert.Equal(t, "foo", proxy.target.Name)
	assert.Equal(t, 0, proxy.target.Age)

	mock2 := proxy.Tap(func(p *proxyMock) {
		p.Name = "bar"
		p.Age = 18
	})
	assert.Equal(t, "bar", proxy.target.Name)
	assert.Equal(t, 18, proxy.target.Age)
	assert.Equal(t, "bar", mock.Name)
	assert.Equal(t, mock, mock2)

	mock3 := proxy.With(func(p *proxyMock) *proxyMock {
		p.Name = "baz"
		p.Age = 20
		return p
	})
	assert.Equal(t, "baz", proxy.target.Name)
	assert.Equal(t, 20, proxy.target.Age)
	assert.Equal(t, "baz", mock.Name)
	assert.Equal(t, mock, mock3)

	mock4 := proxy.When(true, func(p *proxyMock) *proxyMock {
		p.Name = "qux"
		p.Age = 22
		return p
	})
	assert.Equal(t, "qux", proxy.target.Name)
	assert.Equal(t, 22, proxy.target.Age)
	assert.Equal(t, "qux", mock.Name)
	assert.Equal(t, mock, mock4)

	mock5 := proxy.Target()
	assert.Equal(t, mock, mock5)
}

func TestProxy_Struct(t *testing.T) {
	mock := proxyMock{Name: "foo"}

	proxy := NewProxy(mock)
	assert.Equal(t, "foo", proxy.target.Name)
	assert.Equal(t, 0, proxy.target.Age)

	mock2 := proxy.Tap(func(p proxyMock) {
		p.Name = "bar"
		p.Age = 18
	})
	assert.Equal(t, "foo", proxy.target.Name)
	assert.Equal(t, 0, proxy.target.Age)
	assert.Equal(t, "foo", mock2.Name)

	mock3 := proxy.With(func(p proxyMock) proxyMock {
		p.Name = "baz"
		p.Age = 20
		return p
	})
	assert.Equal(t, "baz", proxy.target.Name)
	assert.Equal(t, 20, proxy.target.Age)
	assert.Equal(t, "baz", mock3.Name)
	assert.NotEqual(t, mock, mock3)

	mock4 := proxy.When(true, func(p proxyMock) proxyMock {
		p.Name = "qux"
		p.Age = 22
		return p
	})
	assert.Equal(t, "qux", proxy.target.Name)
	assert.Equal(t, 22, proxy.target.Age)
	assert.Equal(t, "qux", mock4.Name)
	assert.NotEqual(t, mock, mock4)

	mock5 := proxy.Target()
	assert.NotEqual(t, mock, mock5)

	mock6 := proxy.Unless(true, func(p proxyMock) proxyMock {
		p.Name = "quux"
		p.Age = 24
		return p
	})
	assert.Equal(t, "qux", proxy.target.Name)
	assert.Equal(t, 22, proxy.target.Age)
	assert.Equal(t, "qux", mock6.Name)
	assert.NotEqual(t, mock, mock6)
}
