package flatjson

import (
	"encoding/json"
	"strings"
	"testing"
)

type jsonTest struct {
	Name     string
	Input    string
	Expected string
}

var tests = []jsonTest{
	{
		Name:     "empty map",
		Input:    `{}`,
		Expected: `{}`,
	},
	{
		Name:     "empty array",
		Input:    `[]`,
		Expected: `{}`,
	},
	{
		Name:     "one element array",
		Input:    `[true]`,
		Expected: `{"[0]": true}`,
	},
	{
		Name:     "empty map value",
		Input:    `{"foo": {}}`,
		Expected: `{"foo": null}`,
	},
	{
		Name:     "empty array value",
		Input:    `{"foo": []}`,
		Expected: `{"foo": null}`,
	},
	{
		Name:     "array value",
		Input:    `{"foo": [1, 2]}`,
		Expected: `{"foo[0]": 1, "foo[1]": 2}`,
	},
	{
		Name:     "nested map",
		Input:    `{"foo": {"bar": 1}}`,
		Expected: `{"foo.bar": 1}`,
	},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		aux := make(map[string]interface{})

		if err := json.Unmarshal([]byte(test.Expected), &aux); err != nil {
			panic(err)
		}

		size := len(aux)

		r := strings.NewReader(test.Input)

		pairs, err := Parse(r)

		if err != nil {
			t.Error(err)
		}

		if len(pairs) != size {
			t.Errorf("%s: expected %d pairs, got %d", test.Name, size, len(pairs))
			t.Error(pairs)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	var json = `
		{
			"name": "Bob Smith",
			"address": {
				"street": "123 Main Street",
				"city": "Boresville",
				"zipcode": 13943
			},
			"hobbies": ["tennis", "coding", "cooking"]
		}
	`

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(json)
		Parse(r)
	}
}
