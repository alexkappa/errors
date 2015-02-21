package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	if s := err.Error(); s != "test" {
		t.Errorf("expected %q, got %q", "test", s)
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
	// Output: Example failed. write error
}

func ExampleStack() {
	err := New("error with stack trace")
	for _, frame := range err.Stack() {
		fmt.Printf("%s\n", frame.Func)
	}
	// Output: github.com/alexkappa/errors.New
	// github.com/alexkappa/errors.ExampleStack
	// testing.runExample
	// testing.RunExamples
	// testing.(*M).Run
	// main.main
}
