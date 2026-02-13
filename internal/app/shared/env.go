package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var UntrackedEnvVars = []string{"_", "OLDPWD", "SHLVL", "TTY"}

type EnvChange struct {
	Key      string
	OldValue string
	NewValue string
}

type Env struct {
	vars map[string]string
}

func NewEnv(lines []string) *Env {
	vars := make(map[string]string)

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			parts[1] = strings.Trim(parts[1], "\"'")

			vars[parts[0]] = parts[1]
		}
	}

	// remove known vars that are managed by the os/shell and not user apps
	for _, v := range UntrackedEnvVars {
		delete(vars, v)
	}

	return &Env{
		vars: vars,
	}
}

func (old *Env) Diff(new *Env) []EnvChange {
	var changes []EnvChange

	// capture changes when looking up the old key in the new map (changes or removals)
	for oldKey, oldValue := range old.vars {
		newValue, _ := new.vars[oldKey]
		if oldValue != newValue {
			changes = append(changes, EnvChange{Key: oldKey, OldValue: oldValue, NewValue: newValue})
		}
	}

	// check for any new keys that don't exist in the old map (additions)
	for newKey, newValue := range new.vars {
		if oldValue, ok := old.vars[newKey]; !ok {
			changes = append(changes, EnvChange{Key: newKey, OldValue: oldValue, NewValue: newValue})
		}
	}

	return changes
}

func FindLoadPaths(filename string) []string {
	var paths []string

	// currentDir, _ := filepath.Abs(".")
	currentDir, _ := os.Getwd()
	fmt.Println(currentDir)

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

	return paths
}
