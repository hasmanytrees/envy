package zsh

import (
	"bytes"
	"envy/internal/app/shared"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewZsh(t *testing.T) {
	sessionKey := "test-session"
	z := NewZsh(sessionKey)

	if z.SessionKey != sessionKey {
		t.Errorf("expected SessionKey %s, got %s", sessionKey, z.SessionKey)
	}

	homeDir, _ := os.UserHomeDir()
	expectedLoadFilepath := filepath.Join(homeDir, ".cache/envy", "test-session.load.sh")
	expectedUndoFilepath := filepath.Join(homeDir, ".cache/envy", "test-session.undo.sh")

	if z.LoadFilepath != expectedLoadFilepath {
		t.Errorf("expected LoadFilepath %s, got %s", expectedLoadFilepath, z.LoadFilepath)
	}
	if z.UndoFilepath != expectedUndoFilepath {
		t.Errorf("expected UndoFilepath %s, got %s", expectedUndoFilepath, z.UndoFilepath)
	}
}

func TestZsh_Init(t *testing.T) {
	tests := []struct {
		name                  string
		z                     *Zsh
		initScript            string
		checkSessionKey       string
		checkExecLoadFilepath string
		checkExecUndoFilepath string
		checkRmLoadFilepath   string
		checkRmUndoFilepath   string
		wantErr               bool
	}{
		{
			name: "success",
			z: &Zsh{
				SessionKey:   "test-session",
				LoadFilepath: "/tmp/test-session.load.sh",
				UndoFilepath: "/tmp/test-session.undo.sh",
			},
			initScript:            initScript,
			checkSessionKey:       "ENVY_SESSION_KEY=test-session",
			checkExecLoadFilepath: ". /tmp/test-session.load.sh",
			checkExecUndoFilepath: ". /tmp/test-session.undo.sh",
			checkRmLoadFilepath:   "rm /tmp/test-session.load.sh",
			checkRmUndoFilepath:   "rm /tmp/test-session.undo.sh",

			wantErr: false,
		}, {
			name: "error",
			z: &Zsh{
				SessionKey:   "test-session",
				LoadFilepath: "/tmp/test-session.load.sh",
				UndoFilepath: "/tmp/test-session.undo.sh",
			},
			initScript: "{{",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			initScript = tt.initScript

			err := tt.z.Init(&buf)
			if (err != nil) && !tt.wantErr {
				t.Errorf("init error = %v, wantErr %v", err, tt.wantErr)
			}

			output := buf.String()

			if !tt.wantErr {
				if !strings.Contains(output, tt.checkSessionKey) {
					t.Errorf("expected output to contain %q", tt.checkSessionKey)
				}
				if !strings.Contains(output, tt.checkExecLoadFilepath) {
					t.Errorf("expected output to contain %q", tt.checkExecLoadFilepath)
				}
				if !strings.Contains(output, tt.checkExecUndoFilepath) {
					t.Errorf("expected output to contain %q", tt.checkExecUndoFilepath)
				}
				if !strings.Contains(output, tt.checkRmLoadFilepath) {
					t.Errorf("expected output to contain %q", tt.checkRmLoadFilepath)
				}
				if !strings.Contains(output, tt.checkRmUndoFilepath) {
					t.Errorf("expected output to contain %q", tt.checkRmUndoFilepath)
				}
			}
		})
	}
}

func TestZsh_FindLoadPaths(t *testing.T) {
	// this method is simply a wrapper around shared.FindLoadPaths and has
	// no other processing/logic, therefore no real testing is done here
	z := &Zsh{}

	z.FindLoadPaths()
}

func TestZsh_GetSubshellCmd(t *testing.T) {
	z := &Zsh{LoadFilepath: "/tmp/load.sh"}
	cmd := z.GetSubshellCmd()

	expectedArgs := []string{"zsh", "-c", ". /tmp/load.sh; envy export"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d", len(expectedArgs), len(cmd.Args))
	}

	for i := range cmd.Args {
		if cmd.Args[i] != expectedArgs[i] {
			t.Errorf("arg %d: expected %s, got %s", i, expectedArgs[i], cmd.Args[i])
		}
	}
}

func TestZsh_GenLoadFile(t *testing.T) {
	z := &Zsh{LoadFilepath: "/tmp/load.sh"}

	tests := []struct {
		name          string
		paths         []string
		expectedLines []string
	}{
		{
			name:  "no paths",
			paths: []string{},
			expectedLines: []string{
				"#!/bin/zsh",
			},
		},
		{
			name:  "single path",
			paths: []string{"/a/b/envy.sh"},
			expectedLines: []string{
				"#!/bin/zsh",
				". '/a/b/envy.sh'",
			},
		},
		{
			name:  "multiple paths",
			paths: []string{"/a/envy.sh", "/a/b/envy.sh"},
			expectedLines: []string{
				"#!/bin/zsh",
				". '/a/envy.sh'",
				". '/a/b/envy.sh'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, path := z.GenLoadFile(tt.paths)

			if path != z.LoadFilepath {
				t.Errorf("expected path %s, got %s", z.LoadFilepath, path)
			}

			if len(lines) != len(tt.expectedLines) {
				t.Fatalf("expected %d lines, got %d", len(tt.expectedLines), len(lines))
			}

			for i := range lines {
				if lines[i] != tt.expectedLines[i] {
					t.Errorf("line %d: expected %q, got %q", i, tt.expectedLines[i], lines[i])
				}
			}
		})
	}
}

func TestZsh_GenUndoFile(t *testing.T) {
	z := &Zsh{UndoFilepath: "/tmp/undo.sh"}

	tests := []struct {
		name          string
		changes       []shared.EnvChange
		expectedLines []string
	}{
		{
			name:    "no changes",
			changes: []shared.EnvChange{},
			expectedLines: []string{
				"#!/bin/zsh",
			},
		},
		{
			name: "addition (old value empty)",
			changes: []shared.EnvChange{
				{Key: "NEW_VAR", OldValue: "", NewValue: "val"},
			},
			expectedLines: []string{
				"#!/bin/zsh",
				"unset NEW_VAR",
			},
		},
		{
			name: "removal (new value empty)",
			changes: []shared.EnvChange{
				{Key: "OLD_VAR", OldValue: "old_val", NewValue: ""},
			},
			expectedLines: []string{
				"#!/bin/zsh",
				"export OLD_VAR=old_val",
			},
		},
		{
			name: "modification",
			changes: []shared.EnvChange{
				{Key: "MOD_VAR", OldValue: "old", NewValue: "new"},
			},
			expectedLines: []string{
				"#!/bin/zsh",
				"if [[ \"${MOD_VAR}\" == \"new\" ]]; then \n\texport MOD_VAR=old\nfi",
			},
		},
		{
			name: "multiple changes",
			changes: []shared.EnvChange{
				{Key: "ADD", OldValue: "", NewValue: "1"},
				{Key: "REM", OldValue: "2", NewValue: ""},
				{Key: "MOD", OldValue: "3", NewValue: "4"},
			},
			expectedLines: []string{
				"#!/bin/zsh",
				"unset ADD",
				"export REM=2",
				"if [[ \"${MOD}\" == \"4\" ]]; then \n\texport MOD=3\nfi",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, path := z.GenUndoFile(tt.changes)

			if path != z.UndoFilepath {
				t.Errorf("expected path %s, got %s", z.UndoFilepath, path)
			}

			if len(lines) != len(tt.expectedLines) {
				t.Fatalf("expected %d lines, got %d", len(tt.expectedLines), len(lines))
			}

			for i := range lines {
				if lines[i] != tt.expectedLines[i] {
					t.Errorf("line %d: expected %q, got %q", i, tt.expectedLines[i], lines[i])
				}
			}
		})
	}
}
