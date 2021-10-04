package jsonh

import "testing"

func TestEqual(t *testing.T) {
	tests := []struct {
		name       string
		json1      []byte
		json2      []byte
		wantResult bool
	}{
		{
			"blank strings",
			[]byte(""),
			[]byte(""),
			false,
		},
		{
			"invalid first json",
			[]byte(`{"a":1`),
			[]byte(`{"a":1}`),
			false,
		},
		{
			"invalid second json",
			[]byte(`{"a":1}`),
			[]byte(`{"a":1`),
			false,
		},
		{
			"equal strings",
			[]byte(`{"a":1}`),
			[]byte(`{"a":1}`),
			true,
		},
		{
			"same structure with extra spaces",
			[]byte(`{"a":1}`),
			[]byte(`{ "a": 1 }`),
			true,
		},
		{
			"same structure with multiline JSONs",
			[]byte(`{ "a": 1 }`),
			[]byte(`{
					"a": 1
				}`),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Equal(tt.json1, tt.json2)
			if r != tt.wantResult {
				t.Errorf("want to be equal; got not equal")
			}
		})
	}
}
