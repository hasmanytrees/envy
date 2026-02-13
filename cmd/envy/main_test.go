package main

import (
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"envy"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			sout := os.Stdout
			null, _ := os.Open(os.DevNull)
			os.Stdout = null
			defer func() { os.Stdout = sout }()

			// can't easily test main functions that exit through an call to os.Exit
			// therefore this test only executes main without any assertions

			main()
		})
	}
}
