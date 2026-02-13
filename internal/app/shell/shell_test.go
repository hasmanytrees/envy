package shell

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewShell(t *testing.T) {
	tests := []struct {
		name       string
		shellType  string
		sessionKey string
		wantNil    bool
	}{
		{
			name:       "supported shell",
			shellType:  SupportedShellTypes[0],
			sessionKey: "test-session",
			wantNil:    false,
		},
		{
			name:       "unsupported shell",
			shellType:  "unsupported",
			sessionKey: "test-session",
			wantNil:    true,
		},
		{
			name:       "empty shell",
			shellType:  "",
			sessionKey: "test-session",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShell(tt.shellType, tt.sessionKey)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewShell() = %v, wantNil %v", got, tt.wantNil)
			}

			if !tt.wantNil {
				z, ok := got.(*Zsh)
				if !ok {
					t.Errorf("NewShell() did not return a *Zsh for zsh shell type")
				}

				if z.SessionKey != tt.sessionKey {
					t.Errorf("NewShell() SessionKey = %v, want %v", z.SessionKey, tt.sessionKey)
				}
			}
		})
	}
}

func TestFindLoadPaths(t *testing.T) {
	tmp := t.TempDir()
	tmpDir, err := filepath.EvalSymlinks(tmp)
	if err != nil {
		t.Fatal(err)
	}

	// Create directory structure:
	// tmpDir/envy.sh
	// tmpDir/a/ (no file)
	// tmpDir/a/b/envy.sh
	// tmpDir/a/b/c/envy.sh

	dirs := []string{
		filepath.Join(tmpDir, "a"),
		filepath.Join(tmpDir, "a", "b"),
		filepath.Join(tmpDir, "a", "b", "c"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	files := []string{
		filepath.Join(tmpDir, "envy.sh"),
		filepath.Join(tmpDir, "a", "b", "envy.sh"),
		filepath.Join(tmpDir, "a", "b", "c", "envy.sh"),
	}

	for _, file := range files {
		if err := os.WriteFile(file, []byte(""), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", file, err)
		}
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)

	tests := []struct {
		name     string
		workDir  string
		filename string
		want     []string
	}{
		{
			name:     "find in c",
			workDir:  filepath.Join(tmpDir, "a", "b", "c"),
			filename: "envy.sh",
			want: []string{
				filepath.Join(tmpDir, "envy.sh"),
				filepath.Join(tmpDir, "a", "b", "envy.sh"),
				filepath.Join(tmpDir, "a", "b", "c", "envy.sh"),
			},
		},
		{
			name:     "find in b",
			workDir:  filepath.Join(tmpDir, "a", "b"),
			filename: "envy.sh",
			want: []string{
				filepath.Join(tmpDir, "envy.sh"),
				filepath.Join(tmpDir, "a", "b", "envy.sh"),
			},
		},
		{
			name:     "find in a",
			workDir:  filepath.Join(tmpDir, "a"),
			filename: "envy.sh",
			want: []string{
				filepath.Join(tmpDir, "envy.sh"),
			},
		},
		{
			name:     "non-existent file",
			workDir:  filepath.Join(tmpDir, "a", "b", "c"),
			filename: "missing.sh",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Chdir(tt.workDir); err != nil {
				t.Fatalf("failed to chdir to %s: %v", tt.workDir, err)
			}

			got, err := findLoadPaths(tt.filename)
			if err != nil {
				t.Errorf("findLoadPaths() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findLoadPaths() got = %v, want %v", got, tt.want)
			}
		})
	}
}
