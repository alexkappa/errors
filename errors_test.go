package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestBatch(t *testing.T) {
	for batch, expected := range map[BatchError]string{
		NewBatch(New("e1"), New("e2")):                   "multiple errors occured: e1\ne2",
		NewBatch(io.ErrUnexpectedEOF):                    "multiple errors occured: unexpected EOF",
		NewBatch(Wrapf(New("e1"), "m1 %s %s", "b", "z")): "multiple errors occured: m1 b z. e1",
		NewBatch(io.ErrUnexpectedEOF, nil):               "multiple errors occured: unexpected EOF",
		NewBatch():                                       "multiple errors occured: ",
		NewBatch(nil):                                    "multiple errors occured: ",
		NewBatch([]error{nil}...):                        "multiple errors occured: ",
	} {
		if _, ok := batch.(error); !ok {
			t.Error("expected batch to be the generic error")
		}

		if batch.Error() != expected {
			t.Errorf("expected string to be %q, got %s", expected, batch.Error())
		}

		for _, b := range batch.Errors() {
			if s := b.Stack(); s == nil {
				t.Error("empty stack trace")
			}
		}

		if fmt.Sprintf("%s\n", batch) != fmt.Sprintf("%s\n", expected) {
			t.Errorf("expected string to be %q, got %s", fmt.Sprintf("%s\n", expected), fmt.Sprintf("%s\n", batch))
		}

		if fmt.Sprintf("%+s\n", batch) != fmt.Sprintf("%+s\n", expected) {
			t.Errorf("expected string to be %q, got %s", fmt.Sprintf("%+s\n", expected), fmt.Sprintf("%+s\n", batch))
		}

		if fmt.Sprintf("%v", batch) != fmt.Sprintf("%v", expected) {
			t.Errorf("expected string to be %q, got %s", fmt.Sprintf("%v", expected), fmt.Sprintf("%v", batch))
		}
	}
}

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

func TestNew(t *testing.T) {
	for _, message := range []string{
		"foo",
		"bar",
		"baz",
	} {
		if New(message).Error() != message {
			t.Errorf("expected error to be equal to %q", message)
		}
	}
}

func TestErrorf(t *testing.T) {
	for message, args := range map[string][]interface{}{
		"foo %s":       {"f"},
		"bar %s %s %s": {"b", "a", "r"},
		"baz %s %s":    {"b", "z"},
	} {
		expected := fmt.Sprintf(message, args...)

		err := Errorf(message, args...)
		if err.Error() != expected {
			t.Errorf("unexpected error output %s", err.Error())
		}

		expected = fmt.Sprintf("%s. x", expected)

		err = Wrapf(New("x"), message, args...)
		if err.Error() != expected {
			t.Error("unexpected error output %s", err.Error())
		}
	}
}

func TestMarshal(t *testing.T) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(New("test"))
	if err != nil {
		t.Errorf("json error: %s", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"message":"test"`)) {
		t.Errorf("marshaled error doesn't contain expected segment")
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"stack":[]`)) {
		t.Errorf("marshaled error doesn't contain an empty stack trace")
	}
}

func TestMarshalTrace(t *testing.T) {
	MarshalTrace = true
	defer func() {
		MarshalTrace = false
	}()
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(New("test"))
	if err != nil {
		t.Errorf("json error: %s", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"stack":[{"file":"errors.go"`)) {
		t.Errorf("marshaled error should contain a stack trace")
	}
}
