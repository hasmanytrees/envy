package shell

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

var ValidShellTypes = []string{"zsh"}

type Shell interface {
	Init(w io.Writer) error
	FindLoadPaths() ([]string, error)
	GetSubshellCmd() *exec.Cmd
	GenLoadFile(paths []string) ([]string, string)
	GenUndoFile(changes []EnvChange) ([]string, string)
}

func NewShell(shellType string, sessionKey string) Shell {
	switch shellType {
	case "zsh":
		return NewZsh(sessionKey)
	}
	return nil
}

func findLoadPaths(filename string) ([]string, error) {
	var paths []string

	currentDir, _ := filepath.Abs(".")

	for {
		path := filepath.Join(currentDir, filename)
		// check if the file exists ignoring errors for files not found
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
		}

		// get the parent directory
		parentDir := filepath.Dir(currentDir)

		// check if we have reached the root directory (parent is the same as current)
		if parentDir == currentDir {
			break
		}

		// move up to the parent directory for the next iteration
		currentDir = parentDir
	}

	// reverse so that processing can happen naturally (highest directory working down)
	slices.Reverse(paths)

	return paths, nil
}
