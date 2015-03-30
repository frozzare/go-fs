package fs

import (
	"fmt"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

var (
	path, _ = os.Getwd()
	dir     = Open(path)
)

func TestCopy(t *testing.T) {
	err := dir.Copy("files/hello.txt", "files/hello-copy.txt")

	if err != nil {
		panic(err)
	}

	content, err := dir.Read("files/hello-copy.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello, read!\n", content)
}

func TestDelete(t *testing.T) {
	err := dir.Delete("files/hello-copy.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, nil, err)
}

func TestCreateDir(t *testing.T) {
	dir.Delete("files/dir")
	err := dir.CreateDir("files/dir")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, nil, err)
}

func TestGetFileExtension(t *testing.T) {
	assert.Equal(t, "jpg", dir.GetFileExtension("files/284.jpg"))
}

func TestGetSize(t *testing.T) {
	size, err := dir.GetSize("files/284.jpg")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, int64(16361), size)
}

func TestHas(t *testing.T) {
	assert.Equal(t, true, dir.Has("files/284.jpg"))
	assert.Equal(t, false, dir.Has("files/284.png"))
}

func TestListContents(t *testing.T) {
	items := dir.ListContents(".", true)

	fmt.Println(items)

	assert.Equal(t, true, len(items) > 0)
	assert.Equal(t, "fs_test.go", items[len(items)-1].Name)
}

func TestRead(t *testing.T) {
	content, err := dir.Read("files/hello.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello, read!\n", content)

	content, err = dir.Read("files/error")

	assert.Equal(t, "", content)
	assert.NotEqual(t, nil, err)
}

func TestUpdate(t *testing.T) {
	err := dir.Copy("files/hello.txt", "files/hello-update.txt")

	if err != nil {
		panic(err)
	}

	err = dir.Update("files/hello-update.txt", "Hello, update!")

	if err != nil {
		panic(err)
	}

	content, err := dir.Read("files/hello-update.txt")

	assert.Equal(t, "Hello, read!\nHello, update!", content)
	assert.Equal(t, nil, err)
}

func TestWrite(t *testing.T) {
	err := dir.Write("files/hello-write.txt", "Hello, write!\n")

	if err != nil {
		panic(err)
	}

	content, err := dir.Read("files/hello-write.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello, write!\n", content)
	assert.Equal(t, nil, err)
}
