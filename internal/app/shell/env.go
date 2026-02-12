package shell

import (
	"strings"
)

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
	delete(vars, "_")
	delete(vars, "OLDPWD")
	delete(vars, "SHLVL")
	delete(vars, "TTY")

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
