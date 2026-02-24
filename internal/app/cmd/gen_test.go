package cmd

import (
	"context"
	"envy/internal/app/shared"
	"envy/internal/app/shell"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

type fakeShell struct {
	findLoadPaths  func() []string
	getSubshellCmd func() *exec.Cmd
	genLoadFile    func(paths []string) ([]string, string)
	genUndoFile    func(changes []shared.EnvChange) ([]string, string)
}

func (f *fakeShell) Init(_ io.Writer) error {
	return nil
}

func (f *fakeShell) FindLoadPaths() []string {
	return f.findLoadPaths()
}

func (f *fakeShell) GetSubshellCmd() *exec.Cmd {
	return f.getSubshellCmd()
}

func (f *fakeShell) GenLoadFile(paths []string) ([]string, string) {
	return f.genLoadFile(paths)
}

func (f *fakeShell) GenUndoFile(changes []shared.EnvChange) ([]string, string) {
	return f.genUndoFile(changes)
}

func TestGenCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"gen"},
			wantErr: false,
		},
		{
			name:    "with args",
			args:    []string{"gen", "arg1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// using the test shell
			t.Setenv("ENVY_SHELL", "test")
			t.Setenv("ENVY_SESSION_KEY", "12345678")

			rootCmd.SetArgs(tt.args)
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)

			err := genCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("genCmd error = %v, wantErr %v", err, tt.wantErr)
			}

			// cleanup after the test shell
			os.RemoveAll("test.load.sh")
			os.RemoveAll("test.unload.sh")
		})
	}
}

func TestGenPreRun(t *testing.T) {
	tests := []struct {
		name       string
		shellType  string
		sessionKey string
		wantErr    bool
	}{
		{
			name:       "no shell type",
			shellType:  "",
			sessionKey: "12345678",
			wantErr:    true,
		},
		{
			name:       "no session key",
			shellType:  shell.SupportedShellTypes[0],
			sessionKey: "",
			wantErr:    true,
		},
		{
			name:       "no shell type AND no session key",
			shellType:  "",
			sessionKey: "",
			wantErr:    true,
		},
		{
			name:       "with shell type AND session key",
			shellType:  shell.SupportedShellTypes[0],
			sessionKey: "12345678",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ENVY_SHELL", tt.shellType)
			t.Setenv("ENVY_SESSION_KEY", tt.sessionKey)

			err := genPreRun(genCmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("exportRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenRun(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name       string
		fake       *fakeShell
		wantErr    bool
		assertFunc func(t *testing.T)
	}{
		{
			name: "success",
			fake: &fakeShell{
				findLoadPaths:  func() []string { return []string{} },
				getSubshellCmd: func() *exec.Cmd { return exec.Command("sh", "-c", "exit 0") }, // echo full env so diff is empty
				genLoadFile: func(paths []string) ([]string, string) {
					return []string{}, filepath.Join(tmp, "session.load.sh")
				},
				genUndoFile: func(_ []shared.EnvChange) ([]string, string) {
					return []string{}, filepath.Join(tmp, "session.undo.sh")
				},
			},
			wantErr: false,
			assertFunc: func(t *testing.T) {
				// verify files exist
				if _, err := os.Stat(filepath.Join(tmp, "session.load.sh")); err != nil {
					t.Fatalf("load file not created: %v", err)
				}
				if _, err := os.Stat(filepath.Join(tmp, "session.undo.sh")); err != nil {
					t.Fatalf("undo file not created: %v", err)
				}
			},
		},
		{
			name: "error in writeLines due to sh.GenLoadFile",
			fake: &fakeShell{
				findLoadPaths:  func() []string { return []string{} },
				getSubshellCmd: func() *exec.Cmd { return exec.Command("sh", "-c", "exit 0") }, // echo full env so diff is empty
				genLoadFile: func(paths []string) ([]string, string) {
					return []string{}, ""
				},
				genUndoFile: func(_ []shared.EnvChange) ([]string, string) {
					return []string{}, filepath.Join(tmp, "session.undo.sh")
				},
			},
			wantErr: true,
		},
		{
			name: "error subshell.CombinedOutput",
			fake: &fakeShell{
				findLoadPaths:  func() []string { return []string{} },
				getSubshellCmd: func() *exec.Cmd { return exec.Command("sh", "-c", "exit 1") }, // echo full env so diff is empty
				genLoadFile: func(paths []string) ([]string, string) {
					return []string{}, filepath.Join(tmp, "session.load.sh")
				},
				genUndoFile: func(_ []shared.EnvChange) ([]string, string) {
					return []string{}, filepath.Join(tmp, "session.undo.sh")
				},
			},
			wantErr: true,
		},
		{
			name: "error in writeLines due to sh.GenUndoFile",
			fake: &fakeShell{
				findLoadPaths:  func() []string { return []string{} },
				getSubshellCmd: func() *exec.Cmd { return exec.Command("sh", "-c", "exit 0") }, // echo full env so diff is empty
				genLoadFile: func(paths []string) ([]string, string) {
					return []string{}, filepath.Join(tmp, "session.load.sh")
				},
				genUndoFile: func(_ []shared.EnvChange) ([]string, string) {
					return []string{}, ""
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := genCmd.Context()
			ctx = context.WithValue(ctx, "shell", tt.fake)
			genCmd.SetContext(ctx)

			err := genRun(genCmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("genRun() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(t)
			}
		})
	}
}

func TestWriteLines(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "success",
			fileName: path.Join(tmp, "session.load.sh"),
			wantErr:  false,
		},
		{
			name:     "error",
			fileName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeLines([]string{}, tt.fileName)

			if (err != nil) != tt.wantErr {
				t.Errorf("writeLines() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
