package shared

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

func TestNewEnv(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		wantVars map[string]string
	}{
		{
			name:  "simple",
			lines: []string{"FOO=bar", "FIZZ=buzz"},
			wantVars: map[string]string{
				"FOO":  "bar",
				"FIZZ": "buzz",
			},
		},
		{
			name:  "with quotes",
			lines: []string{`FOO="bar"`, `FIZZ='buzz'`, `TEST=var`},
			wantVars: map[string]string{
				"FOO":  "bar",
				"FIZZ": "buzz",
				"TEST": "var",
			},
		},
		{
			name:  "filtered variables",
			lines: []string{"FOO=bar", "_=bash", "OLDPWD=/tmp", "SHLVL=1", "TTY=/dev/pts/0"},
			wantVars: map[string]string{
				"FOO": "bar",
			},
		},
		{
			name:     "no equals",
			lines:    []string{"FOO"},
			wantVars: map[string]string{},
		},
		{
			name:  "multiple equals",
			lines: []string{"FOO=bar=buzz"},
			wantVars: map[string]string{
				"FOO": "bar=buzz",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEnv(tt.lines)
			if !reflect.DeepEqual(got.vars, tt.wantVars) {
				t.Errorf("NewEnv() vars = %v, want %v", got.vars, tt.wantVars)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name string
		old  *Env
		new  *Env
		want []EnvChange
	}{
		{
			name: "no changes",
			old:  &Env{vars: map[string]string{"FOO": "bar"}},
			new:  &Env{vars: map[string]string{"FOO": "bar"}},
			want: nil,
		},
		{
			name: "addition",
			old:  &Env{vars: map[string]string{"FOO": "bar"}},
			new:  &Env{vars: map[string]string{"FOO": "bar", "FIZZ": "buzz"}},
			want: []EnvChange{
				{Key: "FIZZ", OldValue: "", NewValue: "buzz"},
			},
		},
		{
			name: "removal",
			old:  &Env{vars: map[string]string{"FOO": "bar", "FIZZ": "buzz"}},
			new:  &Env{vars: map[string]string{"FOO": "bar"}},
			want: []EnvChange{
				{Key: "FIZZ", OldValue: "buzz", NewValue: ""},
			},
		},
		{
			name: "modification",
			old:  &Env{vars: map[string]string{"FOO": "bar"}},
			new:  &Env{vars: map[string]string{"FOO": "buzz"}},
			want: []EnvChange{
				{Key: "FOO", OldValue: "bar", NewValue: "buzz"},
			},
		},
		{
			name: "multiple changes",
			old:  &Env{vars: map[string]string{"STAY": "same", "CHANGE": "old", "REMOVE": "gone"}},
			new:  &Env{vars: map[string]string{"STAY": "same", "CHANGE": "new", "ADD": "here"}},
			want: []EnvChange{
				{Key: "CHANGE", OldValue: "old", NewValue: "new"},
				{Key: "REMOVE", OldValue: "gone", NewValue: ""},
				{Key: "ADD", OldValue: "", NewValue: "here"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.old.Diff(tt.new)

			// Sort both to compare slices reliably
			sortFunc := func(a, b EnvChange) int {
				if a.Key < b.Key {
					return -1
				}
				if a.Key > b.Key {
					return 1
				}
				return 0
			}
			slices.SortFunc(got, sortFunc)
			slices.SortFunc(tt.want, sortFunc)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
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

			got := FindLoadPaths(tt.filename)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findLoadPaths() got = %v, want %v", got, tt.want)
			}
		})
	}
}
