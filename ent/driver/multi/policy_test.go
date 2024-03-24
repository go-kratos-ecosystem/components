package multi

import (
	"testing"

	"entgo.io/ent/dialect"
	"github.com/stretchr/testify/assert"
)

var driver1, driver2, driver3 dialect.Driver

func TestPolicy_RoundRobinPolicy(t *testing.T) {
	p := RoundRobinPolicy()
	drivers := []dialect.Driver{driver1, driver2, driver3}

	for i := 0; i < 10; i++ {
		assert.Equal(t, drivers[i%3], p.Resolve(drivers))
	}
}

func TestPolicy_RandomPolicy(t *testing.T) {
	p := RandomPolicy()
	drivers := []dialect.Driver{driver1, driver2, driver3}

	for i := 0; i < 10; i++ {
		assert.Contains(t, drivers, p.Resolve(drivers))
	}
}
