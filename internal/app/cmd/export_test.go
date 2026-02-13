package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type errorWriter struct{}

func (e errorWriter) Write(p []byte) (n int, err error) {
	return 0, os.ErrPermission
}

func TestExportCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"export"},
			wantErr: false,
		},
		{
			name:    "with args",
			args:    []string{"export", "arg1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)

			err := exportCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("exportCmd error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExportRun(t *testing.T) {
	key := "TEST_VAR"
	value := "test_value"
	t.Setenv(key, value)

	tests := []struct {
		name    string
		writer  io.Writer
		wantErr bool
	}{
		{
			name:    "success",
			writer:  &bytes.Buffer{},
			wantErr: false,
		},
		{
			name:    "write error",
			writer:  errorWriter{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exportRun(tt.writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("exportRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				buf, ok := tt.writer.(*bytes.Buffer)
				if !ok {
					t.Fatal("writer is not a bytes.Buffer")
				}

				output := buf.String()
				lines := strings.Split(strings.TrimSpace(output), "\n")

				// Check if our test variable is in the output
				found := false
				expectedLine := key + "=" + value
				for _, line := range lines {
					if line == expectedLine {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("expected environment variable %q not found in output", expectedLine)
				}

				// Verify the number of lines matches os.Environ()
				expectedCount := len(os.Environ())
				if len(lines) != expectedCount {
					t.Errorf("expected %d lines, got %d", expectedCount, len(lines))
				}
			}
		})
	}
}
