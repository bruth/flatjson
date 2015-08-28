# flatjson

[![GoDoc](https://godoc.org/github.com/bruth/flatjson?status.svg)](https://godoc.org/github.com/bruth/flatjson)

The `flatjson` package provides an encoder to take a value or existing JSON string and convert it into a flat map or array or key-value pairs. For example, the following document can be flattened to an array of key-value pairs:

```json
{
    "name": "Bob Smith",
    "address": {
        "street": "123 Main Street",
        "city": "Boresville",
        "zipcode": 13943
    },
    "hobbies": ["tennis", "coding", "cooking"]
}
```

```json
[
    ["name", "Bob Smith"],
    ["address.street", "123 Main Street"],
    ["address.city", "Boresville"],
    ["address.zipcode", 13943],
    ["hobbies[0]", "tennis"],
    ["hobbies[1]", "coding"],
    ["hobbies[2]", "cooking"]
]
```

## What is the motivation?

Working with a flat set of a key value pairs is often easier to work with when storing or indexing nested documents. It also makes it easier to compare documents by just iterating over the keys and comparing values rather than recursing into documents or arrays.

## Install

Get the library.

```
go get github.com/bruth/flatjson
```

Install the command line tool.

```
go install github.com/bruth/flatjson/cmd/flatjson
```

## Library Example

Convert existing JSON into flat format.

```go
import (
    "strings"

    "github.com/bruth/flatjson"
)

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

r := strings.NewReader(json)

// Write the JSON to standard out.
if err := flatjson.NewEncoder(os.Stdout).ConvertMap(r); err != nil {
    // Handle error..
}
```

## CLI Tool

```
flatjson [-array] [file]
```

### Example

```
flatjson
{
    "name": "Bob Smith",
    "address": {
        "street": "123 Main Street",
        "city": "Boresville",
        "zipcode": 13943
    },
    "hobbies": ["tennis", "coding", "cooking"]
}
{
    "address.street": "123 Main Street",
    "address.city": "Boresville",
    "address.zipcode": 13943,
    "hobbies[0]": "tennis",
    "hobbies[1]": "coding",
    "hobbies[2]": "cooking",
    "name": "Bob Smith"
    }
