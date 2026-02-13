package shell

import (
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
			name:       "supported shell",
			shellType:  SupportedShellTypes[0],
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShell(tt.shellType, tt.sessionKey)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewShell() = %v, wantNil %v", got, tt.wantNil)
			}

			if !tt.wantNil {
				z, ok := got.(*zsh.Zsh)
				if !ok {
					t.Errorf("NewShell() did not return a *Zsh for zsh shell type")
				}

				if z.SessionKey != tt.sessionKey {
					t.Errorf("NewShell() SessionKey = %v, want %v", z.SessionKey, tt.sessionKey)
				}
			}
		})
	}
}
