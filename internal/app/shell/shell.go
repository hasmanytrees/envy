package shell

import (
	"envy/internal/app/shared"
	"envy/internal/app/shell/zsh"
	"io"
	"os/exec"
)

var SupportedShellTypes = []string{"zsh"}

type Shell interface {
	Init(w io.Writer) error
	FindLoadPaths() []string
	GetSubshellCmd() *exec.Cmd
	GenLoadFile(paths []string) ([]string, string)
	GenUndoFile(changes []shared.EnvChange) ([]string, string)
}

func NewShell(shellType string, sessionKey string) Shell {
	switch shellType {
	case "zsh":
		return zsh.NewZsh(sessionKey)
	}
	return nil
}
