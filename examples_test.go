package errors_test

import (
	"fmt"

	"github.com/alexkappa/errors"
)

func ExampleNew() {
	fmt.Println(errors.New("Example failed."))
	// Output: Example failed. [github.com/alexkappa/errors.New(errors.go:70),github.com/alexkappa/errors_test.ExampleNew(examples_test.go:10),testing.runExample(example.go:99),testing.RunExamples(example.go:36),testing.(*M).Run(testing.go:495),main.main(_testmain.go:66)]
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
	// Output: Example failed. write error [github.com/alexkappa/errors.Wrap(errors.go:82),github.com/alexkappa/errors_test.ExampleWrap(examples_test.go:28),testing.runExample(example.go:99),testing.RunExamples(example.go:36),testing.(*M).Run(testing.go:495),main.main(_testmain.go:66)]
}

func ExampleStack() {
	err := errors.New("error with stack trace")
	for _, frame := range err.Stack() {
		fmt.Printf("%s(%s:%d)\n", frame.Func, frame.File, 0)
	}
	// Output: github.com/alexkappa/errors.New(errors.go:0)
	// github.com/alexkappa/errors_test.ExampleStack(examples_test.go:0)
	// testing.runExample(example.go:0)
	// testing.RunExamples(example.go:0)
	// testing.(*M).Run(testing.go:0)
	// main.main(_testmain.go:0)
}
