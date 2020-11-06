package main

import "testing"

func TestMainSetEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		wantEnv string
	}{
		{"lowercase production", "production", "production"},
		{"lowercase development", "development", "development"},
		{"lowercase test", "test", "test"},
		{"capitalized test", "TEST", "test"},
		{"unknown environment", "unknown", "development"},
	}

	for _, tt := range tests {
		env := setEnv(tt.env)
		if env != tt.wantEnv {
			t.Errorf("want %s; got %s", tt.env, env)
		}
	}
}
