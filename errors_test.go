package errors

import (
	"bytes"
	"encoding/json"
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
