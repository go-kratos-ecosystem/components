package snap

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSnapshot_Basic(t *testing.T) {
	var (
		now   = time.Now()
		value int
	)
	snap := New(func() (int, error) {
		return rand.IntN(1000), nil
	}, Interval[int](time.Millisecond*500))

	for {
		oldValue := value
		value = snap.Get()
		assert.True(t, value >= 0 && value < 1000)

		if oldValue != value {
			break
		}
	}
	assert.True(t, time.Now().Sub(now) < time.Millisecond*1000*2)
}

type User struct {
	Name string
	Age  int
}

func TestSnapshot_Struct(t *testing.T) {
	var (
		internal = time.Millisecond * 500
		now      = time.Now()
		age      int
	)
	snap := New(func() (*User, error) {
		return &User{
			Name: "test",
			Age:  rand.IntN(100),
		}, nil
	}, Interval[*User](internal))

	for {
		oldAge := age
		assert.Equal(t, "test", snap.Get().Name)

		age = snap.Get().Age
		assert.True(t, age >= 0 && age < 100)
		if oldAge != age {
			break
		}
	}
	assert.True(t, time.Now().Sub(now) < internal*2)
}

func TestSnapshot_Error(t *testing.T) {
	snap := New(func() (int, error) {
		return 0, assert.AnError
	})

	assert.Equal(t, 0, snap.Get())
	assert.EqualError(t, snap.Refresh(), assert.AnError.Error())
}
