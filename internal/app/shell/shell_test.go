package shell

import (
	"envy/internal/app/shell/test"
	"envy/internal/app/shell/zsh"
	"testing"
)

func TestNewShell(t *testing.T) {
	tests := []struct {
		name       string
		shellType  string
		sessionKey string
		wantNil    bool
	}{
		{
			name:       "zsh shell",
			shellType:  "zsh",
			sessionKey: "test-session",
			wantNil:    false,
		},
		{
			name:       "unsupported shell",
			shellType:  "unsupported",
			sessionKey: "test-session",
			wantNil:    true,
		},
		{
			name:       "empty shell",
			shellType:  "",
			sessionKey: "test-session",
			wantNil:    true,
		},
		{
			name:       "test shell",
			shellType:  "test",
			sessionKey: "test-session",
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShell(tt.shellType, tt.sessionKey)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewShell() = %v, wantNil %v", got, tt.wantNil)
			}

			if !tt.wantNil {
				if tt.shellType == "zsh" {
					z, ok := got.(*zsh.Zsh)
					if !ok {
						t.Errorf("NewShell() did not return a *Zsh for zsh shell type")
					}

					if z.SessionKey != tt.sessionKey {
						t.Errorf("NewShell() SessionKey = %v, want %v", z.SessionKey, tt.sessionKey)
					}
				} else if tt.shellType == "test" {
					z, ok := got.(*test.Test)
					if !ok {
						t.Errorf("NewShell() did not return a *Test for shell type")
					}

					if z.SessionKey != tt.sessionKey {
						t.Errorf("NewShell() SessionKey = %v, want %v", z.SessionKey, tt.sessionKey)
					}
				}

			}
		})
	}
}
