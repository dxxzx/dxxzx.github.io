// Copyright 2021 magnifier Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/dxxzx/magnifier/core/storage"
	"github.com/dxxzx/magnifier/core/storage/factory"
)

const (
	driverName           = "filesystem"
	defaultRootDirectory = "/var/lib/magnifier"
)

func init() {
	factory.Register(driverName, &filesystemFactory{})
}

type filesystemFactory struct{}

func (f *filesystemFactory) Create(parameters map[string]interface{}) (storage.Driver, error) {
	return FromParameters(parameters)
}

type driver struct {
	root string
}

func FromParameters(parameters map[string]interface{}) (*driver, error) {
	var root = defaultRootDirectory
	if parameters != nil {
		if rootDir, ok := parameters["root"]; ok {
			root = fmt.Sprint(rootDir)
		}
	}
	return &driver{
		root: root,
	}, nil
}

func (d *driver) GetContent(path string) ([]byte, error) {
	rc, err := d.Reader(path, 0)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	p, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *driver) PutContent(path string, content []byte) error {
	writer, err := d.Writer(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, bytes.NewReader(content))
	return err
}

func (d *driver) Reader(subPath string, offset int64) (io.ReadCloser, error) {
	fullPath := d.fullPath(subPath)
	parent := path.Dir(fullPath)

	if _, err := os.Stat(parent); os.IsNotExist(err) {
		os.MkdirAll(parent, 0755)
	}
	file, err := os.OpenFile(fullPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	seekPos, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		file.Close()
		return nil, err
	} else if seekPos < offset {
		file.Close()
		return nil, errors.New("invalid seek offset")
	}

	return file, nil
}

func (d *driver) Writer(subPath string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	fullPath := d.fullPath(subPath)
	parentDir := path.Dir(fullPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(fullPath, flag, perm)
}

func (d *driver) Stat(subPath string) (os.FileInfo, error) {
	fullPath := d.fullPath(subPath)

	fi, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	return fi, nil
}

func (d *driver) Readdir(subPath string) ([]os.FileInfo, error) {
	fullPath := d.fullPath(subPath)

	dir, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	return dir.Readdir(0)
}

func (d *driver) Readdirnames(subPath string) ([]string, error) {
	fullPath := d.fullPath(subPath)

	dir, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	fileNames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(fileNames))
	for _, fileName := range fileNames {
		keys = append(keys, path.Join(subPath, fileName))
	}

	return keys, nil
}

func (d *driver) Move(sourcePath, destPath string) error {
	source := d.fullPath(sourcePath)
	dest := d.fullPath(destPath)

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(path.Dir(dest), 0755); err != nil {
		return err
	}

	err := os.Rename(source, dest)
	return err
}

func (d *driver) Delete(subPath string) error {
	fullPath := d.fullPath(subPath)

	_, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.RemoveAll(fullPath)
	return err
}

func (d *driver) Walk(path string, fn storage.WalkFn) error {
	return factory.Walk(d, path, fn)
}

// fullPath returns the absolute path of a key within the Driver's storage.
func (d *driver) fullPath(subPath string) string {
	return path.Join(d.root, subPath)
}
