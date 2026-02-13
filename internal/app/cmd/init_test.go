package cmd

import (
	"bytes"
	"envy/internal/app/shell"
	"io"
	"testing"
)

func TestInitCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"init"},
			wantErr: true,
		},
		{
			name:    "with 1 valid args",
			args:    []string{"init", shell.SupportedShellTypes[0]},
			wantErr: false,
		},
		{
			name:    "with 2 args",
			args:    []string{"init", "arg1", "arg2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)

			err := initCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("initCmd error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitRun(t *testing.T) {
	tests := []struct {
		name      string
		shellType string
		wantErr   bool
	}{
		{
			name:      "with 1 valid args",
			shellType: shell.SupportedShellTypes[0],
			wantErr:   false,
		},
		{
			name:      "with 1 invalid args",
			shellType: "unsupported-shell",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := initRun(tt.shellType, &bytes.Buffer{})

			if (err != nil) != tt.wantErr {
				t.Errorf("initRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
