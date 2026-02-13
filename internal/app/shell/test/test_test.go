package test

import (
	"bytes"
	"envy/internal/app/shared"
	"testing"
)

func TestNewTest(t *testing.T) {
	test := NewTest()

	if test == nil {
		t.Fatal("NewTest() returned nil")
	}
}

func TestInit(t *testing.T) {
	test := NewTest()

	var buf bytes.Buffer
	err := test.Init(&buf)

	if err != nil {
		t.Errorf("Init() returned unexpected error: %v", err)
	}
}

func TestFindLoadPaths(t *testing.T) {
	test := NewTest()
	paths := test.FindLoadPaths()

	if len(paths) != 0 {
		t.Errorf("FindLoadPaths() = %v, want empty slice", paths)
	}
}

func TestGetSubshellCmd(t *testing.T) {
	test := NewTest()
	cmd := test.GetSubshellCmd()

	if cmd == nil {
		t.Fatal("GetSubshellCmd() returned nil")
	}

	if cmd.Path == "" {
		t.Error("GetSubshellCmd() returned command with empty Path")
	}
}

func TestGenLoadFile(t *testing.T) {
	test := NewTest()

	inputPaths := []string{"/path1", "/path2"}
	paths, filename := test.GenLoadFile(inputPaths)

	if filename != "test.load.sh" {
		t.Errorf("GenLoadFile() filename = %s, want test.load.sh", filename)
	}

	if len(paths) != len(inputPaths) {
		t.Fatalf("GenLoadFile() paths length = %d, want %d", len(paths), len(inputPaths))
	}

	for i := range inputPaths {
		if paths[i] != inputPaths[i] {
			t.Errorf("GenLoadFile() paths[%d] = %s, want %s", i, paths[i], inputPaths[i])
		}
	}
}

func TestGenUndoFile(t *testing.T) {
	test := NewTest()

	changes := []shared.EnvChange{{Key: "FOO", OldValue: "old", NewValue: "new"}}
	paths, filename := test.GenUndoFile(changes)

	if filename != "test.unload.sh" {
		t.Errorf("GenUndoFile() filename = %s, want test.unload.sh", filename)
	}

	if len(paths) != 0 {
		t.Errorf("GenUndoFile() paths = %v, want empty slice", paths)
	}
}
