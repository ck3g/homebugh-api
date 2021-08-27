package main

import (
	"net/url"
	"testing"
)

func TestReadInt(t *testing.T) {
	app := application{}

	tests := []struct {
		name         string
		key          string
		queryString  url.Values
		defaultValue int
		wantValue    int
	}{
		{"successful cast", "key", url.Values{"key": []string{"1"}}, 2, 1},
		{"blank key", "key", url.Values{"key": []string{""}}, 2, 2},
		{"key does not exist", "key", url.Values{}, 2, 2},
		{"cannot convert key to int", "key", url.Values{"key": []string{"string"}}, 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := app.readInt(tt.queryString, tt.key, tt.defaultValue)

			if value != tt.wantValue {
				t.Errorf("want value %d; got %d", tt.wantValue, value)
			}
		})
	}
}
