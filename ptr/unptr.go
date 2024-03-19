package ptr

import "time"

// Deprecated: Use values.Val instead.// Deprecated: Use values.Val instead.
func UnString(s *string) (string, bool) {
	if s == nil {
		return "", false
	}
	return *s, true
}

// Deprecated: Use values.Val instead.
func UnBool(b *bool) (bool, bool) {
	if b == nil {
		return false, false
	}
	return *b, true
}

// Deprecated: Use values.Val instead.
func UnInt(i *int) (int, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

// Deprecated: Use values.Val instead.
func UnInt8(i *int8) (int8, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

// Deprecated: Use values.Val instead.
func UnInt16(i *int16) (int16, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

// Deprecated: Use values.Val instead.
func UnInt32(i *int32) (int32, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

// Deprecated: Use values.Val instead.
func UnInt64(i *int64) (int64, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

// Deprecated: Use values.Val instead.
func UnUint(u *uint) (uint, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

// Deprecated: Use values.Val instead.
func UnUint8(u *uint8) (uint8, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

// Deprecated: Use values.Val instead.
func UnUint16(u *uint16) (uint16, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

// Deprecated: Use values.Val instead.
func UnUint32(u *uint32) (uint32, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

// Deprecated: Use values.Val instead.
func UnUint64(u *uint64) (uint64, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

// Deprecated: Use values.Val instead.
func UnFloat32(f *float32) (float32, bool) {
	if f == nil {
		return 0, false
	}
	return *f, true
}

// Deprecated: Use values.Val instead.
func UnFloat64(f *float64) (float64, bool) {
	if f == nil {
		return 0, false
	}
	return *f, true
}

// Deprecated: Use values.Val instead.
func UnComplex64(c *complex64) (complex64, bool) {
	if c == nil {
		return 0, false
	}
	return *c, true
}

// Deprecated: Use values.Val instead.
func UnComplex128(c *complex128) (complex128, bool) {
	if c == nil {
		return 0, false
	}
	return *c, true
}

// Deprecated: Use values.Val instead.
func UnByte(b *byte) (byte, bool) {
	if b == nil {
		return 0, false
	}
	return *b, true
}

// Deprecated: Use values.Val instead.
func UnRune(r *rune) (rune, bool) {
	if r == nil {
		return 0, false
	}
	return *r, true
}

// Deprecated: Use values.Val instead.
func UnUintptr(u *uintptr) (uintptr, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

// Deprecated: Use values.Val instead.
func UnTime(t *time.Time) (time.Time, bool) {
	if t == nil {
		return time.Time{}, false
	}
	return *t, true
}

// Deprecated: Use values.Val instead.
func UnDuration(d *time.Duration) (time.Duration, bool) {
	if d == nil {
		return 0, false
	}
	return *d, true
}
