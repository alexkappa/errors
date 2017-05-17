package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

// BatchError is an interface that extends the builtin error interface with an
// array of errors.
type BatchError interface {
	error
	Errors() []error
}

// Type batcherrtype is the default implementation of the BatchError interface.
// It is not exported so users can only use it via the NewBatch function.
type batcherrtype struct {
	errors []error
}

// Error is an interface that extends the builtin error interface with inner
// errors and stack traces.
type Error interface {
	// Message returns the error message of the error.
	Message() string
	// Inner returns the inner error that this error wraps.
	Inner() error
	// Stack returns the stack trace that led to the error.
	Stack() Frames

	error
}

var (
	// This setting enables a stack trace to be printed when an Error is being
	// marshaled.
	//
	// 	err := errors.New("example")
	// 	b, _ := json.Marshal(err) // {"message":"example","stack":[{"file":"<file>","line":<line>,"func": "<function>"},...]}
	MarshalTrace = false
)

// Type errtype is the default implementation of the Error interface. It is not
// exported so users can only use it via the New or Wrap functions.
type errtype struct {
	message string
	inner   error
	stack   Frames
}

// Message returns the error message of the error.
func (t *errtype) Message() string {
	return t.message
}

// Inner returns the inner error that this error wraps.
func (t *errtype) Stack() Frames {
	return t.stack
}

// Stack returns the stack trace that led to the error.
func (t *errtype) Inner() error {
	return t.inner
}

// Error implements the standard library error interface.
func (t *errtype) Error() string {
	if t.inner != nil {
		return t.message + ". " + t.inner.Error()
	}
	return t.message
}

// Format implements the standard library fmt.Formatter interface. Credit to
// Dave Cheney's github.com/pkg/errors.
func (t *errtype) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "%s", t.message)
		if t.inner != nil {
			fmt.Fprintf(s, ". %s", t.inner.Error())
		}
		if s.Flag('+') {
			fmt.Fprint(s, "\n")
			fmt.Fprintf(s, "%s", t.stack.String())
		}
	case 's':
		fmt.Fprintf(s, "%s", t.message)
		if t.inner != nil {
			fmt.Fprintf(s, ". %s", t.inner.Error())
		}
	case 'q':
		fmt.Fprintf(s, "%q", t.message)
	}
}

func (t *errtype) MarshalJSON() ([]byte, error) {
	b, err := t.stack.MarshalJSON()
	if err != nil {
		return b, err
	}

	var buf bytes.Buffer
	fmt.Fprint(&buf, "{")
	fmt.Fprintf(&buf, `"message":%q`, t.message)
	if t.inner != nil {
		fmt.Fprintf(&buf, `,"inner":%q`, t.inner)
	}
	if t.stack != nil {
		fmt.Fprintf(&buf, `,"stack":%s`, b)
	}
	fmt.Fprint(&buf, "}")

	return buf.Bytes(), nil
}

func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// NewBatch creates a new BatchError.
func NewBatch(errs []error) BatchError {
	return &batcherrtype{
		errors: errs,
	}
}

//Errors returns the errors
func (b *batcherrtype) Errors() []error {
	return b.errors
}

// Error implements the standard library error interface.
func (b *batcherrtype) Error() string {
	var msg string

	size := len(b.Errors())
	if size > 0 {
		msg = flatten(b.Errors())
	}

	return msg
}

func flatten(e []error) string {
	if size := len(e); size > 1 {
		return e[0].Error() + "\n" + flatten(e[1:])
	}

	return e[0].Error()
}

// New creates a new Error with the supplied message.
func New(message string) Error {
	return new(message, 3)
}

// Errorf creates a new Error with the supplied message and arguments formatted
// in the manner of fmt.Printf.
func Errorf(message string, args ...interface{}) Error {
	return new(sprintf(message, args...), 3)
}

func new(message string, skip int) Error {
	return &errtype{
		message: message,
		stack:   Stack(skip),
	}
}

// Wrap creates a new Error that wraps err.
func Wrap(err error, message string) Error {
	return wrap(err, message, 3)
}

// Wrapf creates a new Error that wraps err formatted in the manner of
// fmt.Printf.
func Wrapf(err error, message string, args ...interface{}) Error {
	return wrap(err, sprintf(message, args...), 3)
}

func wrap(err error, message string, skip int) Error {
	if errT, ok := err.(*errtype); ok {
		errT.stack = nil // drop the stack trace of the inner error.
	} else {
		err = &errtype{message: err.Error()}
	}
	return &errtype{
		message: message,
		inner:   err,
		stack:   Stack(skip),
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
		buf.WriteByte('\t')
		buf.WriteString(frame.Func)
		buf.WriteByte('(')
		buf.WriteString(frame.File)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))
		buf.WriteByte(')')
		if i < len(f)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func (f Frames) MarshalJSON() ([]byte, error) {
	if !MarshalTrace {
		return []byte("[]"), nil
	}
	b := bytes.NewBuffer(nil)
	b.WriteByte('[')
	e := json.NewEncoder(b)
	for i, frame := range f {
		err := e.Encode(frame)
		if err != nil {
			return b.Bytes(), err
		}
		if i+1 < len(f) {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.Bytes(), nil
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
			File: filepath.Base(file),
			Line: line,
			Func: fn.Name(),
		}
	}
	return frames
}
