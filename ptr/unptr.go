package ptr

import "time"

func UnString(s *string) (string, bool) {
	if s == nil {
		return "", false
	}
	return *s, true
}

func UnBool(b *bool) (bool, bool) {
	if b == nil {
		return false, false
	}
	return *b, true
}

func UnInt(i *int) (int, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

func UnInt8(i *int8) (int8, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

func UnInt16(i *int16) (int16, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

func UnInt32(i *int32) (int32, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

func UnInt64(i *int64) (int64, bool) {
	if i == nil {
		return 0, false
	}
	return *i, true
}

func UnUint(u *uint) (uint, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

func UnUint8(u *uint8) (uint8, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

func UnUint16(u *uint16) (uint16, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

func UnUint32(u *uint32) (uint32, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

func UnUint64(u *uint64) (uint64, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

func UnFloat32(f *float32) (float32, bool) {
	if f == nil {
		return 0, false
	}
	return *f, true
}

func UnFloat64(f *float64) (float64, bool) {
	if f == nil {
		return 0, false
	}
	return *f, true
}

func UnComplex64(c *complex64) (complex64, bool) {
	if c == nil {
		return 0, false
	}
	return *c, true
}

func UnComplex128(c *complex128) (complex128, bool) {
	if c == nil {
		return 0, false
	}
	return *c, true
}

func UnByte(b *byte) (byte, bool) {
	if b == nil {
		return 0, false
	}
	return *b, true
}

func UnRune(r *rune) (rune, bool) {
	if r == nil {
		return 0, false
	}
	return *r, true
}

func UnUintptr(u *uintptr) (uintptr, bool) {
	if u == nil {
		return 0, false
	}
	return *u, true
}

func UnTime(t *time.Time) (time.Time, bool) {
	if t == nil {
		return time.Time{}, false
	}
	return *t, true
}

func UnDuration(d *time.Duration) (time.Duration, bool) {
	if d == nil {
		return 0, false
	}
	return *d, true
}
