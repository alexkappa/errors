package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	err := New("test")

	if m := err.Message(); m != "test" {
		t.Errorf("expected %q, got %q", "test", m)
	}

	if i := err.Inner(); i != nil {
		t.Errorf("unexpected inner error %q", i)
	}

	if s := err.Stack(); s == nil {
		t.Error("empty stack trace")
	}

	if s := err.Error(); strings.Index(s, "test") != 0 {
		t.Errorf("expected string to start with %q", "test")
	}
}

func TestJSON(t *testing.T) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(New("test"))
	if err != nil {
		t.Errorf("json error: %s", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"message":"test"`)) {
		t.Errorf("marshaled error doesn't contain expected segment")
	}
}

func ExampleNew() {
	fmt.Println(New("Example failed."))
	// Output: Example failed. [github.com/alexkappa/errors.New(errors.go:64),github.com/alexkappa/errors.ExampleNew(errors_test.go:43),testing.runExample(example.go:99),testing.RunExamples(example.go:36),testing.(*M).Run(testing.go:486),main.main(_testmain.go:60)]
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
		err = Wrap(err, "Example failed")
	}
	fmt.Println(err)
	// Output: Example failed. write error [github.com/alexkappa/errors.Wrap(errors.go:76),github.com/alexkappa/errors.ExampleWrap(errors_test.go:61),testing.runExample(example.go:99),testing.RunExamples(example.go:36),testing.(*M).Run(testing.go:486),main.main(_testmain.go:60)]
}

func ExampleStack() {
	err := New("error with stack trace")
	for _, frame := range err.Stack() {
		fmt.Printf("%s(%s:%d)\n", frame.Func, frame.File, 0)
	}
	// Output: github.com/alexkappa/errors.New(errors.go:0)
	// github.com/alexkappa/errors.ExampleStack(errors_test.go:0)
	// testing.runExample(example.go:0)
	// testing.RunExamples(example.go:0)
	// testing.(*M).Run(testing.go:0)
	// main.main(_testmain.go:0)
}
