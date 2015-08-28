// The flatjson package supplies types for converting nested JSON structures
// into flat representations.
//
// For example, the following document can be flattened to an array of key-value pairs:
//
//    {
//        "name": "Bob Smith",
//        "address": {
//            "street": "123 Main Street",
//            "city": "Boresville",
//            "zipcode": 13943
//        },
//        "hobbies": ["tennis", "coding", "cooking"]
//    }
//
// will be flattened to:
//
//    [
//        ["name", "Bob Smith"],
//        ["address.street", "123 Main Street"],
//        ["address.city", "Boresville"],
//        ["address.zipcode", 13943],
//        ["hobbies.[0]", "tennis"],
//        ["hobbies.[1]", "coding"],
//        ["hobbies.[2]", "cooking"]
//    ]
//

package flatjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Pair is a key-value Pair of JSON tokens.
type Pair struct {
	Key   string
	Value interface{}
}

func (p *Pair) String() string {
	return fmt.Sprintf("[%s: %v]", p.Key, p.Value)
}

type tokArray [2]interface{}

type arrayPairs []*Pair

func (a arrayPairs) MarshalJSON() ([]byte, error) {
	aux := make([]tokArray, len(a))

	for i, p := range a {
		aux[i] = tokArray{json.Token(p.Key), p.Value}
	}

	return json.Marshal(aux)
}

// Pairs is a set of key-value pairs.
type mapPairs []*Pair

func (m mapPairs) MarshalJSON() ([]byte, error) {
	aux := make(map[string]json.Token, len(m))

	for _, p := range m {
		aux[p.Key] = p.Value
	}

	return json.Marshal(aux)
}

// JSON delimiters.
var (
	lbrace  = json.Delim('{')
	rbrace  = json.Delim('}')
	lsquare = json.Delim('[')
	rsquare = json.Delim(']')

	pathd = "."
)

// parseJSON decodes a JSON-encoded value into a set of pairs.
func parseJSON(r io.Reader) ([]*Pair, error) {
	var (
		// Current token.
		tok json.Token

		// Current converted key-value Pair.
		key   string
		value interface{}

		// Set of key-value pairs.
		pairs []*Pair

		err error

		// Denotes the decoder just entered an map or array.
		inmap bool
		inarr bool

		// Denotes whether the current map or array is empty.
		empty bool

		// Denotes the next token will be an map key.
		onkey bool

		// The current index in the array.
		arridx = []byte("[0]")
		arrkey string

		// Pre-allocate 10 levels deep
		path = make([]string, 10)
		dest []string

		pos = -1
	)

	dec := json.NewDecoder(r)

	for {
		tok, err = dec.Token()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// Evaluate the token to determine next key-value pair.
		switch tok.(type) {
		case json.Delim:
			switch tok {
			case lbrace:
				empty = true
				inmap = true
				onkey = true
				pos++

				// Double the size
				if pos == len(path) {
					dest = make([]string, pos*2)
					copy(dest, path)
					path = dest
				}

			case rbrace:
				if empty && pos > 0 {
					pairs = append(pairs, &Pair{
						Key: strings.Join(path[:pos+1], pathd),
					})
				}

				inmap = false
				// This is here because the map may be empty.
				onkey = true
				pos--

			case lsquare:
				empty = true
				inarr = true

				// Reset the array index.
				arridx[1] = '0'

				if pos > 0 {
					arrkey = path[pos]
				} else {
					arrkey = ""
				}

			case rsquare:
				if empty && pos >= 0 {
					pairs = append(pairs, &Pair{
						Key: strings.Join(path[:pos+1], pathd),
					})
				}

				inarr = false
			}

		// Keys and values.
		default:
			empty = false

			// The current token is the key of a map
			if onkey {
				// Add to key path and increment the position.
				path[pos] = tok.(string)
				onkey = false

				// Token is a map or array value.
			} else {
				value = tok

				if inarr {
					// Only occurs when the top-level value is an array.
					if pos < 0 {
						pos = 0
					}

					path[pos] = arrkey + string(arridx)
					arridx[1]++
				} else if inmap {
					onkey = true
				}

				// Serialize path into key.
				key = strings.Join(path[:pos+1], pathd)

				pairs = append(pairs, &Pair{
					Key:   key,
					Value: value,
				})
			}
		}
	}

	return pairs, nil
}

// Encoder encodes a value into a flat JSON map or array.
type Encoder struct {
	w io.Writer
}

// EncodeArray encodes a value as a flat JSON array.
func (f *Encoder) EncodeArray(v interface{}) error {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	pairs, err := parseJSON(buf)

	if err != nil {
		return err
	}

	return json.NewEncoder(f.w).Encode(arrayPairs(pairs))
}

// EncodeMap encodes a value as a flat JSON map.
func (f *Encoder) EncodeMap(v interface{}) error {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	pairs, err := parseJSON(buf)

	if err != nil {
		return err
	}

	return json.NewEncoder(f.w).Encode(mapPairs(pairs))
}

// ConvertArray re-encodes a JSON value into a flat array.
func (f *Encoder) ConvertArray(r io.Reader) error {
	pairs, err := parseJSON(r)

	if err != nil {
		return err
	}

	return json.NewEncoder(f.w).Encode(arrayPairs(pairs))
}

// ConvertMap re-encodes a JSON value into a flat map.
func (f *Encoder) ConvertMap(r io.Reader) error {
	pairs, err := parseJSON(r)

	if err != nil {
		return err
	}

	return json.NewEncoder(f.w).Encode(mapPairs(pairs))
}

// NewEncoder initializes a new Encoder for the writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// EncodeMap encodes a value into a flat JSON map.
func EncodeMap(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)

	if err := enc.EncodeMap(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// EncodeArray encodes a value into a flat JSON array.
func EncodeArray(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)

	if err := enc.EncodeArray(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ConvertMap re-encodes JSON into a flat map.
func ConvertMap(r io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)

	if err := enc.ConvertMap(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ConvertArray re-encodes JSON into a flat array.
func ConvertArray(r io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)

	if err := enc.ConvertArray(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Parse returns a slice of key-value pairs.
func Parse(r io.Reader) ([]*Pair, error) {
	return parseJSON(r)
}
