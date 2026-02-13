package shell

import (
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
