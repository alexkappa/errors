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
	// 	github.com/alexkappa/errors.New(errors.go:164)
	// 	github.com/alexkappa/errors_test.ExamplePrintfvplus(examples_test.go:20)
	// 	testing.runExample(example.go:115)
	// 	testing.RunExamples(example.go:38)
	// 	testing.(*M).Run(testing.go:744)
	// 	main.main(_testmain.go:78)
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
	// testing.RunExamples(example.go:0)
	// testing.(*M).Run(testing.go:0)
	// main.main(_testmain.go:0)
}

func ExampleBatchStack() {
	batch := errors.NewBatch(errors.New("error with stack trace"))
	for _, b := range batch.Errors() {
		for _, frame := range b.Stack() {
			fmt.Printf("%s(%s:%d)\n", frame.Func, frame.File, 0)
		}
	}
	// Output: github.com/alexkappa/errors.Wrap(errors.go:0)
	//github.com/alexkappa/errors.NewBatch(errors.go:0)
	//github.com/alexkappa/errors_test.ExampleBatchStack(examples_test.go:0)
	//testing.runExample(example.go:0)
	//testing.RunExamples(example.go:0)
	//testing.(*M).Run(testing.go:0)
	//main.main(_testmain.go:0)
}
