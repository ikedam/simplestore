package simplestore

import "fmt"

// ProgrammingError indicates an error caused by specifying inappropriate value
type ProgrammingError struct {
	msg string
}

// NewProgrammingError returns a new ProgirammingError
func NewProgrammingError(msg string) *ProgrammingError {
	return &ProgrammingError{
		msg: msg,
	}
}

// NewProgrammingErrorf returns a new ProgirammingError
func NewProgrammingErrorf(format string, a ...any) *ProgrammingError {
	return &ProgrammingError{
		msg: fmt.Sprintf(format, a...),
	}
}

// Error is an implementation for error
func (e *ProgrammingError) Error() string {
	return e.msg
}
