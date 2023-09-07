package ptr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	assert.Equal(t, true, *Bool(true))
	assert.Equal(t, false, *Bool(false))
}

func TestInt(t *testing.T) {
	assert.Equal(t, 1, *Int(1))
}

func TestInt8(t *testing.T) {
	assert.Equal(t, int8(1), *Int8(1))
}

func TestInt16(t *testing.T) {
	assert.Equal(t, int16(1), *Int16(1))
}

func TestInt32(t *testing.T) {
	assert.Equal(t, int32(1), *Int32(1))
}

func TestInt64(t *testing.T) {
	assert.Equal(t, int64(1), *Int64(1))
}

func TestUint(t *testing.T) {
	assert.Equal(t, uint(1), *Uint(1))
}

func TestUint8(t *testing.T) {
	assert.Equal(t, uint8(1), *Uint8(1))
}

func TestUint16(t *testing.T) {
	assert.Equal(t, uint16(1), *Uint16(1))
}

func TestUint32(t *testing.T) {
	assert.Equal(t, uint32(1), *Uint32(1))
}

func TestUint64(t *testing.T) {
	assert.Equal(t, uint64(1), *Uint64(1))
}

func TestFloat32(t *testing.T) {
	assert.Equal(t, float32(1), *Float32(1))
}

func TestFloat64(t *testing.T) {
	assert.Equal(t, float64(1), *Float64(1))
}

func TestComplex64(t *testing.T) {
	assert.Equal(t, complex64(1), *Complex64(1))
}

func TestComplex128(t *testing.T) {
	assert.Equal(t, complex128(1), *Complex128(1))
}

func TestTime(t *testing.T) {
	now := time.Now()
	assert.Equal(t, now, *Time(now))
}

func TestDuration(t *testing.T) {
	assert.Equal(t, time.Second, *Duration(time.Second))
}

func TestByte(t *testing.T) {
	assert.Equal(t, byte(1), *Byte(byte(1)))
}

func TestRune(t *testing.T) {
	assert.Equal(t, rune(1), *Rune(rune(1)))
}

func TestUintptr(t *testing.T) {
	assert.Equal(t, uintptr(1), *Uintptr(uintptr(1)))
}

func TestString(t *testing.T) {
	assert.Equal(t, "test", *String("test"))
}
