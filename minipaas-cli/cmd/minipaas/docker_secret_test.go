package main

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestSecretCreate_NameAndLiteralCall(t *testing.T) {
	called := false
	var gotName string
	var gotBytes []byte
	dockerSecretInspect = func(name string) error { return assertErr }
	dockerSecretCreate = func(name string, content []byte, verbose bool) error {
		called = true
		gotName = name
		gotBytes = append([]byte(nil), content...)
		return nil
	}
	t.Cleanup(func() {
		dockerSecretInspect = func(n string) error { return nil }
		dockerSecretCreate = func(name string, content []byte, verbose bool) error { return nil }
	})

	content := []byte("supersecret")
	name, err := secretCreate("env", content, false)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	h := sha256.Sum256(content)
	prefix := hex.EncodeToString(h[:])[:8]
	want := "env." + prefix
	if name != want {
		t.Fatalf("name mismatch: %q vs %q", name, want)
	}
	if !called || gotName != want || !reflect.DeepEqual(gotBytes, content) {
		t.Fatalf("secretCreate did not forward correctly")
	}
}

func TestSecretCreateLiteral_ExistsSkips(t *testing.T) {
	invoked := false
	dockerSecretInspect = func(name string) error { return nil }
	dockerSecretCreate = func(name string, content []byte, verbose bool) error { invoked = true; return nil }
	t.Cleanup(func() {
		dockerSecretInspect = func(n string) error { return nil }
		dockerSecretCreate = func(name string, content []byte, verbose bool) error { return nil }
	})

	if err := secretCreateLiteral("sec", []byte("x"), false); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if invoked {
		t.Fatalf("runner should not be invoked when exists=true")
	}
}
