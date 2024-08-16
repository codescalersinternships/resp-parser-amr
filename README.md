# RESP Parser

The RESP Parser is a Go library that parses the Redis Serialization Protocol (RESP) format. It allows you to decode various RESP data types including strings, errors, integers, bulk strings, and arrays.

## Installation
To install the package, use `go get`:

```bash
go get github.com/codescalersinternships/resp-parser-amr

```

## Usage

### Import the RESP Parser package

```go
import resp "github.com/codescalersinternships/resp-parser-amr/pkg"
```

### Reading RESP Data
Use the Read function to read and parse RESP data from an io.Reader. The function returns a Value representing the parsed data and an error if something goes wrong.
```go
r := resp.NewResp(<yourReader>)
value, err := r.Read()
if err != nil {
    // Handle error
}
// Process the parsed value

```
### Supported RESP Data Types

- String
- Error
- Integer
- Bulk String
- Array

## Example 
Here's a simple example of how to use the RESP parser:
``` go
package main

import (
    "bytes"
    "fmt"
    "github.com/codescalersinternships/resp-parser-amr/pkg"
)

func main() {
    input := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
    r := resp.NewResp(bytes.NewBufferString(input))
    
    value, err := r.Read()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Parsed value: %+v\n", value)
}

```

## Testing
To run the tests for this package, use the following command:

```bash
make test
```
