package operator

import (
	"io"
)

// Operator is a common interface for a flow control.
type Operator interface {
	io.Closer
	Take() error
	Give()
}
