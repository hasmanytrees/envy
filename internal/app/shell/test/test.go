package test

import (
	"envy/internal/app/shared"
	"io"
	"os/exec"
)

type Test struct{}

func NewTest() *Test {
	return &Test{}
}

func (t *Test) Init(w io.Writer) error {
	return nil
}

func (t *Test) FindLoadPaths() []string {
	return []string{}
}

func (t *Test) GetSubshellCmd() *exec.Cmd {
	return exec.Command("echo", "testing 1, 2, 3")
}

func (t *Test) GenLoadFile(paths []string) ([]string, string) {
	return paths, "test.load.sh"
}

func (t *Test) GenUndoFile(changes []shared.EnvChange) ([]string, string) {
	return []string{}, "test.unload.sh"
}
