package fs

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestCopy(t *testing.T) {
	err := Copy("test/files/hello.txt", "test/files/hello-copy.txt")

	if err != nil {
		panic(err)
	}

	content, err := Read("test/files/hello-copy.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello, read!\n", content)
}

func TestDelete(t *testing.T) {
	err := Delete("test/files/hello-copy.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, nil, err)
}

func TestCreateDir(t *testing.T) {
	Delete("test/files/dir")
	err := CreateDir("test/files/dir")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, nil, err)
}

func TestGetFileExtension(t *testing.T) {
	assert.Equal(t, "jpg", GetFileExtension("test/files/284.jpg"))
}

func TestGetSize(t *testing.T) {
	size, err := GetSize("test/files/284.jpg")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, int64(16361), size)
}

func TestExists(t *testing.T) {
	assert.Equal(t, nil, Exists("test/files/284.jpg"))
	assert.NotEqual(t, nil, Exists("test/files/284.png"))
	assert.Equal(t, nil, Exists("test/dir"))
}

func TestListContents(t *testing.T) {
	items := ListContents(".", true)

	assert.Equal(t, true, len(items) > 0)
	assert.Equal(t, "hello.txt", items[len(items)-1].Name)
}

func TestRead(t *testing.T) {
	content, err := Read("test/files/hello.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello, read!\n", content)

	content, err = Read("test/files/error")

	assert.Equal(t, "", content)
	assert.NotEqual(t, nil, err)
}

func TestUpdate(t *testing.T) {
	err := Copy("test/files/hello.txt", "test/files/hello-update.txt")

	if err != nil {
		panic(err)
	}

	err = Update("test/files/hello-update.txt", "Hello, update!")

	if err != nil {
		panic(err)
	}

	content, err := Read("test/files/hello-update.txt")

	assert.Equal(t, "Hello, read!\nHello, update!", content)
	assert.Equal(t, nil, err)
}

func TestWrite(t *testing.T) {
	err := Write("test/files/hello-write.txt", "Hello, write!\n")

	if err != nil {
		panic(err)
	}

	content, err := Read("test/files/hello-write.txt")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Hello, write!\n", content)
	assert.Equal(t, nil, err)
}
