package main

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfigDsn(t *testing.T) {
	f, err := os.Open("./testdata/config.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		environment string
		wantDSN     string
	}{
		{"production", "root:super-secure@/homebugh?parseTime=true"},
		{"development", "root:@/homebugh_development?parseTime=true"},
		{"test", "root:@/homebugh_test?parseTime=true"},
	}

	for _, tt := range tests {
		t.Run(tt.environment, func(t *testing.T) {
			d := cfg.dsn(tt.environment)
			if d != tt.wantDSN {
				t.Errorf("want DSN %s; got %s", tt.wantDSN, d)
			}
		})
	}
}
