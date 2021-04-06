package parallel

import (
	"errors"
)

type ErrTimeoutExceeded struct{}

func (e *ErrTimeoutExceeded) Error() string {
	return "timeout exceeded"
}

func IsTimeoutExceeded(err error) bool {
	return errors.Is(err, &ErrTimeoutExceeded{})
}
