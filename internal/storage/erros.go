package storage

import (
	"errors"
	"fmt"
)

type URLError struct {
	Err error
}

type URLGoneError struct {
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

func (e *URLGoneError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func NewURLGoneError(errMsg string) error {
	err := errors.New(errMsg)
	return &URLGoneError{
		Err: err,
	}
}
