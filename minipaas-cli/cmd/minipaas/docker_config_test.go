package main

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestConfigCreate_NameAndLiteralCall(t *testing.T) {
	called := false
	var gotName string
	var gotBytes []byte
	dockerConfigInspect = func(name string) error { return assertErr } // signal not exists
	dockerConfigCreate = func(name string, content []byte, verbose bool) error {
		called = true
		gotName = name
		gotBytes = append([]byte(nil), content...)
		return nil
	}
	t.Cleanup(func() {
		dockerConfigInspect = func(n string) error { return nil }
		dockerConfigCreate = func(name string, content []byte, verbose bool) error { return nil }
	})

	content := []byte("hello")
	name, err := configCreate("app.conf", content, false)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	h := sha256.Sum256(content)
	prefix := hex.EncodeToString(h[:])[:8]
	want := "app.conf." + prefix
	if name != want {
		t.Fatalf("name mismatch: got %q want %q", name, want)
	}
	if !called || gotName != want || !reflect.DeepEqual(gotBytes, content) {
		t.Fatalf("configCreate did not forward correctly")
	}
}

var assertErr = &fakeError{}

type fakeError struct{}

func (e *fakeError) Error() string { return "assert" }

func TestConfigCreateLiteral_ExistsSkips(t *testing.T) {
	invoked := false
	dockerConfigInspect = func(name string) error { return nil } // exists
	dockerConfigCreate = func(name string, content []byte, verbose bool) error { invoked = true; return nil }
	t.Cleanup(func() {
		dockerConfigInspect = func(n string) error { return nil }
		dockerConfigCreate = func(name string, content []byte, verbose bool) error { return nil }
	})

	if err := configCreateLiteral("cfg", []byte("x"), false); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if invoked {
		t.Fatalf("runner should not be invoked when exists=true")
	}
}
