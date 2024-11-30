# utf8reader

A simple go package that converts an io.Reader to a utf8 encoded io.Reader.
It automatically detects the encoding of the input and converts it to utf8.

## Usage

```go

package main

import (
    "fmt"
    "bytes"

    "github.com/kpym/utf8reader"
)

func main() {
    // Create a reader with koi8-r encoded "Това е на български"
    r := bytes.NewReader([]byte{0xF4, 0xCF, 0xD7, 0xC1, 0x20, 0xC5, 0x20, 0xCE, 0xC1, 0x20, 0xC2, 0xDF, 0xCC, 0xC7, 0xC1, 0xD2, 0xD3, 0xCB, 0xC9})
    reader := utf8reader.New(r)

    // Read the content of the reader
    buf := make([]byte, 100)
    n, err := reader.Read(buf)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(buf[:n]))
    // Output: Това е на български
}
```

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/kpym/utf8reader.svg)](https://pkg.go.dev/github.com/kpym/utf8reader)

You can find the documentation on [pkg.go.dev](https://pkg.go.dev/github.com/kpym/utf8reader).

## License

[MIT](LICENSE)

