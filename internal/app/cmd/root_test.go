package cmd

import (
	"io"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "success no args",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "error bad args",
			args:    []string{"arg1", "arg2", "arg3"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)

			err := Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("exportCmd error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
