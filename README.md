# go-github-markdown [![Build Status](https://travis-ci.org/frozzare/go-fs.svg?branch=master)](https://travis-ci.org/frozzare/go-fs)

 Work easy with files on the local filesystem.
 Current path will be appended to the given path if the path don't start with `/`.

 Not tested on Windows.

 View the [docs](http://godoc.org/github.com/frozzare/go-fs).

## Installation

```
$ go get github.com/frozzare/go-fs
```

## Example

```go
package main

import (
    "log"

	"github.com/frozzare/go-fs"
)

func main() {
    err := fs.Write("files/hello.txt", "Hello, world!")

    if err != nil {
        log.Fatal(err)
    }

    content, err := fs.Read("files/hello.txt")

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(content)
    // Hello, world!\n
}
```

# License

MIT © [Fredrik Forsmo](https://github.com/frozzare)
