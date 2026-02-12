package shell

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

//go:embed zsh_init.zsh
var initScript string

type Zsh struct {
	SessionKey   string
	LoadFilepath string
	UndoFilepath string
}

func NewZsh(sessionKey string) *Zsh {
	homeDir, _ := os.UserHomeDir()

	loadFilepath := filepath.Join(homeDir, ".cache/envy", fmt.Sprintf("%s.load.sh", sessionKey))
	undoFilepath := filepath.Join(homeDir, ".cache/envy", fmt.Sprintf("%s.undo.sh", sessionKey))

	return &Zsh{
		SessionKey:   sessionKey,
		LoadFilepath: loadFilepath,
		UndoFilepath: undoFilepath,
	}
}

func (z *Zsh) Init(w io.Writer) error {
	t, err := template.New("init").Parse(initScript)
	if err != nil {
		return err
	}

	return t.Execute(w, z)
}

func (z *Zsh) FindLoadPaths() ([]string, error) {
	return findLoadPaths("envy.sh")
}

func (z *Zsh) GetSubshellCmd() *exec.Cmd {
	return exec.Command("zsh", "-c", fmt.Sprintf(". %s; envy export", z.LoadFilepath))
}

func (z *Zsh) GenLoadFile(paths []string) ([]string, string) {
	var lines []string

	lines = append(lines, "#!/bin/zsh")

	for _, path := range paths {
		lines = append(lines, fmt.Sprintf(". '%s'", path))
	}

	return lines, z.LoadFilepath
}

func (z *Zsh) GenUndoFile(changes []EnvChange) ([]string, string) {
	var lines []string

	lines = append(lines, "#!/bin/zsh")

	for _, change := range changes {
		if len(change.OldValue) == 0 {
			// this was an addition so we need to remove it now
			// unset FOO
			lines = append(lines, fmt.Sprintf("unset '%s'", change.Key))
		} else if len(change.NewValue) == 0 {
			// this was a removal so we need to add it back now
			// export FOO=bar
			lines = append(lines, fmt.Sprintf("export '%s=%s'", change.Key, change.OldValue))
		} else {
			// this was a change so we need to change it back now as long as it hasn't been changed outside of envy (thus the check)
			// if [[ $FOO == "baz" ]]; then
			//	export FOO=bar
			// fi
			lines = append(lines, fmt.Sprintf("if [[ \"${%s}\" == \"%s\" ]]; then \n\texport %s=%s\nfi", change.Key, change.NewValue, change.Key, change.OldValue))
		}
	}

	return lines, z.UndoFilepath
}
