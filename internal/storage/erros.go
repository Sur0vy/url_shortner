package storage

import (
	"errors"
	"fmt"
)

type URLError struct {
	Err error
}

func (e *URLError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewURLError(errMsg string) error {
	err := errors.New(errMsg)
	return &URLError{
		Err: err,
	}
}
