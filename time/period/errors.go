package period

import "errors"

// ErrInvalidRange indicates that an interval's end precedes its start.
var ErrInvalidRange = errors.New("period: end must not precede start")
