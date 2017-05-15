package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestBatchError(t *testing.T) {
	batch := NewBatch()

	errors := []Error{New("test1"), New("test2")}

	batch.Append(New("test1"))
	batch.Append(New("test2"))

	for i, r := range batch.Errors() {
		if errors[i].Error() != r.Error() {
			t.Errorf("expected error %q to be equal", "test")
		}
	}

	if batch.IsEmpty() {
		t.Error("expected batch to not be empty")
	}

	if _, ok := batch.(error); !ok {
		t.Error("expected batch to be the generic error")
	}

	s := batch.Error()
	if strings.Compare(s, "test1;test2") != 0 {
		t.Errorf("expected string to start with %q, got %s", "test1;test2", s)
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
