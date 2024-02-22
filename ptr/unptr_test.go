package ptr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnString(t *testing.T) {
	s1, ok1 := UnString(String("test"))
	assert.True(t, ok1)
	assert.Equal(t, "test", s1)

	s2, ok2 := UnString(nil)
	assert.False(t, ok2)
	assert.Equal(t, "", s2)
}

func TestUnInt(t *testing.T) {
	i1, ok1 := UnInt(Int(1))
	assert.True(t, ok1)
	assert.Equal(t, 1, i1)

	i2, ok2 := UnInt(nil)
	assert.False(t, ok2)
	assert.Equal(t, 0, i2)
}

func TestUnInt8(t *testing.T) {
	i1, ok1 := UnInt8(Int8(1))
	assert.True(t, ok1)
	assert.Equal(t, int8(1), i1)

	i2, ok2 := UnInt8(nil)
	assert.False(t, ok2)
	assert.Equal(t, int8(0), i2)
}

func TestUnInt16(t *testing.T) {
	i1, ok1 := UnInt16(Int16(1))
	assert.True(t, ok1)
	assert.Equal(t, int16(1), i1)

	i2, ok2 := UnInt16(nil)
	assert.False(t, ok2)
	assert.Equal(t, int16(0), i2)
}

func TestUnInt32(t *testing.T) {
	i1, ok1 := UnInt32(Int32(1))
	assert.True(t, ok1)
	assert.Equal(t, int32(1), i1)

	i2, ok2 := UnInt32(nil)
	assert.False(t, ok2)
	assert.Equal(t, int32(0), i2)
}

func TestUnInt64(t *testing.T) {
	i1, ok1 := UnInt64(Int64(1))
	assert.True(t, ok1)
	assert.Equal(t, int64(1), i1)

	i2, ok2 := UnInt64(nil)
	assert.False(t, ok2)
	assert.Equal(t, int64(0), i2)
}

func TestUnUint(t *testing.T) {
	i1, ok1 := UnUint(Uint(1))
	assert.True(t, ok1)
	assert.Equal(t, uint(1), i1)

	i2, ok2 := UnUint(nil)
	assert.False(t, ok2)
	assert.Equal(t, uint(0), i2)
}

func TestUnUint8(t *testing.T) {
	i1, ok1 := UnUint8(Uint8(1))
	assert.True(t, ok1)
	assert.Equal(t, uint8(1), i1)

	i2, ok2 := UnUint8(nil)
	assert.False(t, ok2)
	assert.Equal(t, uint8(0), i2)
}

func TestUnUint16(t *testing.T) {
	i1, ok1 := UnUint16(Uint16(1))
	assert.True(t, ok1)
	assert.Equal(t, uint16(1), i1)

	i2, ok2 := UnUint16(nil)
	assert.False(t, ok2)
	assert.Equal(t, uint16(0), i2)
}

func TestUnUint32(t *testing.T) {
	i1, ok1 := UnUint32(Uint32(1))
	assert.True(t, ok1)
	assert.Equal(t, uint32(1), i1)

	i2, ok2 := UnUint32(nil)
	assert.False(t, ok2)
	assert.Equal(t, uint32(0), i2)
}

func TestUnUint64(t *testing.T) {
	i1, ok1 := UnUint64(Uint64(1))
	assert.True(t, ok1)
	assert.Equal(t, uint64(1), i1)

	i2, ok2 := UnUint64(nil)
	assert.False(t, ok2)
	assert.Equal(t, uint64(0), i2)
}

func TestUnFloat32(t *testing.T) {
	f1, ok1 := UnFloat32(Float32(1.1))
	assert.True(t, ok1)
	assert.Equal(t, float32(1.1), f1)

	f2, ok2 := UnFloat32(nil)
	assert.False(t, ok2)
	assert.Equal(t, float32(0), f2)
}

func TestUnFloat64(t *testing.T) {
	f1, ok1 := UnFloat64(Float64(1.1))
	assert.True(t, ok1)
	assert.Equal(t, 1.1, f1)

	f2, ok2 := UnFloat64(nil)
	assert.False(t, ok2)
	assert.Equal(t, 0.0, f2)
}

func TestUnBool(t *testing.T) {
	b1, ok1 := UnBool(Bool(true))
	assert.True(t, ok1)
	assert.Equal(t, true, b1)

	b2, ok2 := UnBool(nil)
	assert.False(t, ok2)
	assert.Equal(t, false, b2)
}

func TestUnComplex64(t *testing.T) {
	c1, ok1 := UnComplex64(Complex64(1 + 1i))
	assert.True(t, ok1)
	assert.Equal(t, complex64(1+1i), c1)

	c2, ok2 := UnComplex64(nil)
	assert.False(t, ok2)
	assert.Equal(t, complex64(0), c2)
}

func TestUnComplex128(t *testing.T) {
	c1, ok1 := UnComplex128(Complex128(1 + 1i))
	assert.True(t, ok1)
	assert.Equal(t, 1+1i, c1)

	c2, ok2 := UnComplex128(nil)
	assert.False(t, ok2)
	assert.Equal(t, 0+0i, c2)
}

func TestUnByte(t *testing.T) {
	b1, ok1 := UnByte(Byte(1))
	assert.True(t, ok1)
	assert.Equal(t, byte(1), b1)

	b2, ok2 := UnByte(nil)
	assert.False(t, ok2)
	assert.Equal(t, byte(0), b2)
}

func TestUnRune(t *testing.T) {
	r1, ok1 := UnRune(Rune('a'))
	assert.True(t, ok1)
	assert.Equal(t, 'a', r1)

	r2, ok2 := UnRune(nil)
	assert.False(t, ok2)
	assert.Equal(t, rune(0), r2)
}

func TestUnUintptr(t *testing.T) {
	u1, ok1 := UnUintptr(Uintptr(1))
	assert.True(t, ok1)
	assert.Equal(t, uintptr(1), u1)

	u2, ok2 := UnUintptr(nil)
	assert.False(t, ok2)
	assert.Equal(t, uintptr(0), u2)
}

func TestUnTime(t *testing.T) {
	now := time.Now()
	t1, ok1 := UnTime(Time(now))
	assert.True(t, ok1)
	assert.Equal(t, now, t1)

	t2, ok2 := UnTime(nil)
	assert.False(t, ok2)
	assert.Equal(t, time.Time{}, t2)
}

func TestUnDuration(t *testing.T) {
	d1, ok1 := UnDuration(Duration(time.Second))
	assert.True(t, ok1)
	assert.Equal(t, time.Second, d1)

	d2, ok2 := UnDuration(nil)
	assert.False(t, ok2)
	assert.Equal(t, time.Duration(0), d2)
}
