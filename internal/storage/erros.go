package storage

type URLExError struct {
	message string
}

type URLGoneError struct {
	message string
}

func (e *URLExError) Error() string {
	return e.message
}

func NewURLExError() *URLExError {
	return &URLExError{
		message: "URL is exist",
	}
}

func (e *URLGoneError) Error() string {
	return e.message
}

func NewURLGoneError() *URLGoneError {
	return &URLGoneError{
		message: "URL is gone",
	}
}
