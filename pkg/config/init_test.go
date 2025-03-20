package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name       string
		envValue   string
		wantPort   int64
		wantErr    bool
		setupFunc  func()
		teardownFn func()
	}{
		{
			name:     "empty environment variable",
			envValue: "",
			wantPort: 8080,
			setupFunc: func() {
				os.Setenv(PORT_ENV, "")
			},
			teardownFn: func() {
				os.Unsetenv(PORT_ENV)
			},
		},
		{
			name:     "valid port number",
			envValue: "12345",
			wantPort: 12345,
			setupFunc: func() {
				os.Setenv(PORT_ENV, "12345")
			},
			teardownFn: func() {
				os.Unsetenv(PORT_ENV)
			},
		},
		{
			name:     "invalid port number",
			envValue: "abc",
			wantPort: 8080,
			setupFunc: func() {
				os.Setenv(PORT_ENV, "abc")
			},
			teardownFn: func() {
				os.Unsetenv(PORT_ENV)
			},
		},
		{
			name:     "unset environment variable",
			envValue: "",
			wantPort: 8080,
			setupFunc: func() {
				os.Unsetenv(PORT_ENV)
			},
			teardownFn: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setupFunc()

			// Run the test
			config := NewConfig()
			if config.Port != tt.wantPort {
				t.Errorf("NewConfig().Port = %d, want %d", config.Port, tt.wantPort)
			}

			// Teardown
			tt.teardownFn()
		})
	}
}
