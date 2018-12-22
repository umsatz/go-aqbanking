package aqbanking

import "fmt"

// Error represents an error message with a code
type Error struct {
	Message string
	Code    int
}

func newError(message string, code _Ctype_int) *Error {
	return &Error{message, int(code)}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %d", e.Message, e.Code)
}
