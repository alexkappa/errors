package errors_test

import (
	"fmt"

	"github.com/alexkappa/errors"
)

func ExamplePrintfv() {
	fmt.Printf("%s", errors.New("Example failed."))
	// Output: Example failed.
}

func ExamplePrintfs() {
	fmt.Printf("%v", errors.New("Example failed."))
	// Output: Example failed.
}

func ExamplePrintfvplus() {
	fmt.Printf("%+v", errors.New("Example failed."))
	// Output: Example failed.
	// 	github.com/alexkappa/errors.New(errors.go:114)
	// 	github.com/alexkappa/errors_test.ExamplePrintfvplus(examples_test.go:20)
	// 	testing.runExample(example.go:123)
	// 	testing.runExamples(example.go:46)
	// 	testing.(*M).Run(testing.go:823)
	// 	main.main(_testmain.go:62)
}

type errWriter uint8

func (e errWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}

var (
	w errWriter
	p []byte
)

func ExampleWrap() {
	_, err := w.Write(p)
	if err != nil {
		err = errors.Wrap(err, "Example failed")
	}
	fmt.Println(err)
	// Output: Example failed. write error
}

func ExampleStack() {
	err := errors.New("error with stack trace")
	for _, frame := range err.Stack() {
		fmt.Printf("%s(%s:%d)\n", frame.Func, frame.File, 0)
	}
	// Output: github.com/alexkappa/errors.New(errors.go:0)
	// github.com/alexkappa/errors_test.ExampleStack(examples_test.go:0)
	// testing.runExample(example.go:0)
	// testing.runExamples(example.go:0)
	// testing.(*M).Run(testing.go:0)
	// main.main(_testmain.go:0)
}
