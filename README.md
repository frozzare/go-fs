# go-fs [![Build Status](https://travis-ci.org/frozzare/go-fs.svg?branch=master)](https://travis-ci.org/frozzare/go-fs)

 Work easy with files and directories on the local filesystem.

 Not tested on Windows.

 View the [docs](http://godoc.org/github.com/frozzare/go-fs).

## Installation

```
$ go get github.com/frozzare/go-fs
```

## Example

Current path will be appended to the given path if the path don't start with `/`.

```go
package main

import (
    "log"

    "github.com/frozzare/go-fs"
)

func main() {
    err := fs.Write("test/files/example.txt", "Hello, world!")

    if err != nil {
        log.Fatal(err)
    }

    content, err := fs.Read("test/files/example.txt")

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(content)
    // Hello, world!
}
```

# License

MIT © [Fredrik Forsmo](https://github.com/frozzare)
