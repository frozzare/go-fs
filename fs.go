package fs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// ContentItem is used for list of contents
type ContentItem struct {
	Name string
	Path string
	Type string
}

// Get the real path of the file or directory
func realPath(file string) string {
	if string(file[0]) == "/" {
		return file
	}

	if string(file[0]) != "/" {
		file = "/" + file
	}

	_, filename, _, _ := runtime.Caller(3)
	dir := path.Join(path.Dir(filename), file)

	if _, err := os.Stat(dir); err == nil && strings.HasSuffix(dir, file) {
		return dir
	}

	current, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	dir = file

	if strings.HasSuffix(dir, current) {
		return dir
	}

	return current + dir
}

// Copy will return true if success otherwise a error
//
// - `src` the source file to copy from
// - `dest` the target file to copy to.
func Copy(src, dest string) error {
	srcFile, err := os.Open(realPath(src))

	if err != nil {
		return err
	}

	defer srcFile.Close()

	sfi, err := srcFile.Stat()

	if err != nil {
		return err
	}

	if !sfi.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	destFile, err := os.Create(realPath(dest))

	if err != nil {
		return err
	}

	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

// CreateDir will create a directory
//
// - `dir` the directory path to create
// - `args` file permission uint32 (optional)
func CreateDir(dir string, args ...interface{}) error {
	permission := uint32(0644)

	if len(args) > 0 {
		permission = args[0].(uint32)
	}

	dir = realPath(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, os.FileMode(permission))
}

// Delete will delete the file or directory
//
// - `path` the path to delete
func Delete(path string) error {
	return os.Remove(realPath(path))
}

// GetFileExtension will return the file extension
//
// - `file` the file path to get file extension from
func GetFileExtension(file string) string {
	ext := filepath.Ext(realPath(file))

	if ext == "" {
		return ""
	}

	return ext[1:]
}

// GetSize returns the size of the file or error
//
// - `file` the file path to get file size from
func GetSize(file string) (int64, error) {
	stat, err := os.Stat(realPath(file))

	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

// Exists returns nil when file exists and error when file does not exist.
//
// - `path` the path to check if exists or not
func Exists(path string) error {
	_, err := os.Stat(realPath(path))

	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return err
	}

	return nil
}

// ListContents will return a list of contents (files and directories)
//
// Args:
// - path to directory that should be listed (optional, default "./*")
// - recursive bool (optional, default false)
func ListContents(args ...interface{}) []ContentItem {
	dir := "./*"
	path := "/"
	recursive := false
	files := []string{}

	if string(dir[0]) != "." {
		path = path + dir
	}

	if len(args) > 0 {
		if len(args[0].(string)) > 0 {
			dir = realPath(args[0].(string)) + "/*"
		}

		recursive = len(args) > 1 && args[1] != nil
	}

	var err error

	if recursive {
		dir = dir[0 : len(dir)-2]
		err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if f.Name() == "." || f.Name() == ".." {
				return nil
			}

			files = append(files, path)

			return nil
		})
	} else {
		files, err = filepath.Glob(dir)
	}

	if err != nil {
		log.Fatal(err)
	}

	result := make([]ContentItem, len(files))

	for i, s := range files {
		fi, err := os.Stat(s)

		if err != nil {
			continue
		}

		pt := "File"

		if fi.IsDir() {
			pt = "Directory"
		}

		if string(s[0]) == "/" {
			s = s[1:]
		}

		item := ContentItem{
			Name: fi.Name(),
			Path: realPath(path + s),
			Type: pt,
		}

		result[i] = item
	}

	return result
}

// Read will return the content of the file or error
//
// - `file` is the file path to read from
func Read(file string) (string, error) {
	err := Exists(file)

	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadFile(realPath(file))

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ReadJSON into a interface{} of your object
//
// - `file` the file path to read JSON from
// - `as` the interface to read to
func ReadJSON(file string, as interface{}) error {
	err := Exists(file)

	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(realPath(file))

	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &as); err != nil {
		return err
	}

	return nil
}

// Update will append text to file
//
// - `file` the file path to append content to
// - `body` the string to write
// - `args` fiel permission uint32 (optional)
func Update(file, content string, args ...interface{}) error {
	permission := uint32(0644)

	if len(args) > 0 {
		permission = args[0].(uint32)
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.FileMode(permission))

	if err != nil {
		return err
	}

	if _, err := f.WriteString(content); err != nil {
		return err
	}

	return nil
}

// Write will write text to file
//
// - `file` is the file path to write to
// - `body` is the string to write
// - `args` fiel permission uint32 (optional)
func Write(file, content string, args ...interface{}) error {
	permission := uint32(0644)

	if len(args) > 0 {
		permission = args[0].(uint32)
	}

	err := ioutil.WriteFile(realPath(file), []byte(content), os.FileMode(permission))

	if err != nil {
		return err
	}

	return nil
}

// WriteJSON will write JSON of the `content` interface.
//
// - `file` is the file path to write to
// - `content` is the interface of your JSON
// - `args` file permission uint32 (optional)
func WriteJSON(file string, content interface{}, args ...interface{}) error {
	body, err := json.Marshal(content)

	if err != nil {
		return err
	}

	return Write(file, string(body), args...)
}
