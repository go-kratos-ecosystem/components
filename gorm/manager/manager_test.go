package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	m                  *Manager
	db1, db2, db3, db4 *gorm.DB
)

func init() {
	m = New(db1)
	m.Register("db2", db2)
	m.Register("db3", db3)
	m.Register("db4", db4)
}

func TestManager(t *testing.T) {
	assert.Equal(t, db1, m.DB)
	assert.Equal(t, db1, m.Conn())
	assert.Equal(t, db2, m.Conn("db2"))
	assert.Equal(t, db3, m.Conn("db3"))
	assert.Equal(t, db4, m.Conn("db4"))
	assert.Panics(t, func() {
		m.Conn("db5")
	})
}

func BenchmarkManager(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = m.DB
			_ = m.Conn()
			_ = m.Conn("db2")
			_ = m.Conn("db3")
			_ = m.Conn("db4")
		}
	})
}
