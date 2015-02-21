package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

type Error interface {
	// Message returns the error message of the error.
	Message() string
	// Inner returns the inner error that this error wraps.
	Inner() error
	// Stack returns the stack trace that led to the error.
	Stack() Frames
	// Error satisfies the standard library error interface.
	Error() string
}

// Type T is the default implementation of the Error interface. Users shouldn't
// need to use this struct directly.
type T struct {
	M string `json:"message"`
	I error  `json:"inner,omitempty"`
	S Frames `json:"stack,omitempty"`
}

// Message returns the error message of the error.
func (t *T) Message() string {
	return t.M
}

// Inner returns the inner error that this error wraps.
func (t *T) Stack() Frames {
	return t.S
}

// Stack returns the stack trace that led to the error.
func (t *T) Inner() error {
	return t.I
}

// Error satisfies the standard library error interface.
func (t *T) Error() string {
	if t.I != nil {
		return fmt.Sprintf("%s. %s", t.M, t.I)
	}
	return fmt.Sprintf("%s", t.M)
}

// New creates a new Error with the supplied message.
func New(message string) Error {
	return new(message, 3)
}

// Newf creates a new Error with the supplied formating.
func Newf(format string, v ...interface{}) Error {
	return new(fmt.Sprintf(format, v...), 3)
}

func new(message string, skip int) Error {
	return &T{
		M: message,
		S: Stack(skip),
	}
}

// Wrap creates a new Error that wraps err.
func Wrap(err error, message string) Error {
	return wrap(err, message, 3)
}

// Wrapf creates a new Error that wraps err.
func Wrapf(err error, format string, v ...interface{}) Error {
	return wrap(err, fmt.Sprintf(format, v...), 3)
}

func wrap(err error, message string, skip int) Error {
	if errT, ok := err.(*T); ok {
		errT.S = nil // drop the stack trace of the inner error.
	} else {
		err = &T{M: err.Error()}
	}
	return &T{
		M: message,
		I: err,
		S: Stack(skip),
	}
}

// Frame contains information for a single stack frame.
type Frame struct {
	// File is the path to the file of the caller.
	File string `json:"file"`
	// Line is the line in the file where the function call was made.
	Line int `json:"line"`
	// Func is the name of the caller.
	Func string `json:"func"`
}

type Frames []Frame

// String is used to satisfy the fmt.Stringer interface. It formats the stack
// trace as a comma separated list of "file:line function".
func (f Frames) String() string {
	var buf bytes.Buffer
	for i, frame := range f {
		buf.WriteString(fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Func))
		if i < len(f)-1 {
			buf.WriteByte(',')
		}
	}
	return buf.String()
}

// Stack returns the stack trace of the function call, while skipping the first
// skip frames.
func Stack(skip int) Frames {
	callers := make([]uintptr, 10)
	n := runtime.Callers(skip, callers)
	callers = callers[:n-2] // skip runtime.main and runtime.goexit function calls
	frames := make(Frames, len(callers))
	for i, caller := range callers {
		fn := runtime.FuncForPC(caller)
		file, line := fn.FileLine(caller)
		frames[i] = Frame{
			File: file,
			Line: line,
			Func: fn.Name(),
		}
	}
	return frames
}
