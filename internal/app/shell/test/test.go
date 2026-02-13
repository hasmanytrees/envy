package test

import (
	"envy/internal/app/shared"
	"io"
	"os/exec"
)

type Test struct {
	SessionKey string
}

func NewTest(sessionKey string) *Test {
	return &Test{
		SessionKey: sessionKey,
	}
}

func (t *Test) Init(_ io.Writer) error {
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

func (t *Test) GenUndoFile(_ []shared.EnvChange) ([]string, string) {
	return []string{}, "test.unload.sh"
}
