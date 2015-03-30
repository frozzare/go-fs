package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Directory struct contains the path
type Directory struct {
	Path string
}

// ContentItem is used for list of contents
type ContentItem struct {
	Name string
	Path string
	Type string
}

// Get the real path of the file or directory
func (d *Directory) realPath(path string) string {
	if string(path[0]) != "/" {
		path = "/" + path
	}

	if strings.Contains(path, d.Path) {
		return path
	}

	return d.Path + path
}

// Copy will return true if success otherwise a error
func (d *Directory) Copy(src, dest string) error {
	srcFile, err := os.Open(d.realPath(src))

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

	destFile, err := os.Create(d.realPath(dest))

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
func (d *Directory) CreateDir(dir string, args ...interface{}) error {
	permission := uint32(0644)

	if len(args) > 0 {
		permission = args[0].(uint32)
	}

	dir = d.realPath(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, os.FileMode(permission))
}

// Delete will delete the file or directory
func (d *Directory) Delete(path string) error {
	return os.Remove(d.realPath(path))
}

// GetFileExtension will return the file extension
func (d *Directory) GetFileExtension(file string) string {
	ext := filepath.Ext(d.realPath(file))

	if ext == "" {
		return ""
	}

	return ext[1:]
}

// GetSize returns the size of the file or error
func (d *Directory) GetSize(file string) (int64, error) {
	stat, err := os.Stat(d.realPath(file))

	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

// Has returns true or false if the file exists in the directory or not.
func (d *Directory) Has(file string) bool {
	_, err := os.Stat(d.realPath(file))
	return err == nil
}

// ListContents will return a list of contents (files and directories)
func (d *Directory) ListContents(args ...interface{}) []ContentItem {
	dir := "./*"
	path := "/"
	recursive := false
	files := []string{}

	if string(dir[0]) != "." {
		path = path + dir
	}

	if len(args) > 0 {
		if len(args[0].(string)) > 0 {
			dir = d.realPath(args[0].(string)) + "/*"
		}

		recursive = len(args) > 1 && args[1] != nil
	}

	var err error

	if recursive {
		dir = dir[0 : len(dir)-2]
		fmt.Println(dir)
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
			Path: d.realPath(path + s),
			Type: pt,
		}

		result[i] = item
	}

	return result
}

// Open will return a new instance of Directory struct
func Open(path string) *Directory {
	return &Directory{path}
}

// Read will return the content of the file or error
func (d *Directory) Read(file string) (string, error) {
	if d.Has(file) == false {
		return "", fmt.Errorf("The file %s doesn't exists", file)
	}

	content, err := ioutil.ReadFile(d.realPath(file))

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Update will append text to file
func (d *Directory) Update(file, content string, args ...interface{}) error {
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
func (d *Directory) Write(file, content string, args ...interface{}) error {
	permission := uint32(0644)

	if len(args) > 0 {
		permission = args[0].(uint32)
	}

	err := ioutil.WriteFile(d.realPath(file), []byte(content), os.FileMode(permission))

	if err != nil {
		return err
	}

	return nil
}
