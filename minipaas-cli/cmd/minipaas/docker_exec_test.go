package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestDockerContainerExec_BuildsArgs(t *testing.T) {
	var got []string
	_runCommand = func(cmd []string, verbose bool) error {
		got = append([]string{}, cmd...)
		return nil
	}
	t.Cleanup(func() { _runCommand = runCommand })

	err := dockerContainerExec("abc123", []string{"sh", "-c", "echo ok"}, false)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []string{"docker", "exec", "-i", "abc123", "sh", "-c", "echo ok"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch\n got:%#v\nwant:%#v", got, want)
	}
}

func TestDockerContainerExecOutput_BuildsArgs(t *testing.T) {
	var got []string
	_runCommandOutput = func(cmd []string, verbose bool) (string, error) {
		got = append([]string{}, cmd...)
		return "output", nil
	}
	t.Cleanup(func() { _runCommandOutput = runCommandOutput })

	_, err := dockerContainerExecOutput("cid", []string{"env"}, true)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []string{"docker", "exec", "-i", "cid", "env"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch\n got:%#v\nwant:%#v", got, want)
	}
}

func TestGetContainerID(t *testing.T) {
	// success: returns first non-empty id
	_runCommandOutput = func(cmd []string, verbose bool) (string, error) {
		return "id1\nid2\n", nil
	}
	t.Cleanup(func() { _runCommandOutput = runCommandOutput })

	id, err := getContainerID("svc")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != "id1" {
		t.Fatalf("id mismatch: %q", id)
	}

	// no ids
	_runCommandOutput = func(cmd []string, verbose bool) (string, error) {
		return "\n", nil
	}
	_, err = getContainerID("svc")
	if err == nil {
		t.Fatalf("expected error when no containers")
	}

	// command failure
	_runCommandOutput = func(cmd []string, verbose bool) (string, error) {
		return "", errors.New("boom")
	}
	_, err = getContainerID("svc")
	if err == nil {
		t.Fatalf("expected error on failure")
	}
}
