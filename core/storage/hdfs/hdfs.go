// Copyright 2021 magnifier Author. All Rights Reserved.
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
package hdfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/colinmarc/hdfs/v2"
	"github.com/dxxzx/magnifier/core/storage"
	"github.com/dxxzx/magnifier/core/storage/factory"
)

const (
	driverName           = "hdfs"
	defaultRootDirectory = "/user/magnifier"
	defaultUser          = "magnifier"
)

func init() {
	factory.Register(driverName, hdfsFactory{})
}

type hdfsFactory struct{}

func (h hdfsFactory) Create(parameters map[string]interface{}) (storage.Driver, error) {
	return FromParameters(parameters)
}

type driver struct {
	rootDirectory string
	client        *hdfs.Client
}

func FromParameters(parameters map[string]interface{}) (*driver, error) {
	var (
		err           error
		rootDirectory = defaultRootDirectory
		user          = defaultUser
		namenodes     []string
	)

	if parameters != nil {
		if rootDir, ok := parameters["root"]; ok {
			rootDirectory = fmt.Sprint(rootDir)
		}
		if nodes, ok := parameters["namenodes"]; ok {
			switch v := nodes.(type) {
			case []interface{}:
				for _, node := range v {
					namenodes = append(namenodes, node.(string))
				}
			default:
				return nil, fmt.Errorf("namenode config error %v", v)
			}
		}
		if u, ok := parameters["user"]; ok {
			user = fmt.Sprint(u)
		}
	}

	options := hdfs.ClientOptions{
		Addresses: namenodes,
		User:      user,
	}
	client, err := hdfs.NewClient(options)
	if err != nil {
		return nil, err
	}
	driver := &driver{
		rootDirectory: rootDirectory,
		client:        client,
	}

	return driver, nil
}

func (d *driver) GetContent(subPath string) ([]byte, error) {
	realPath := d.fullPath(subPath)
	file, err := d.client.Open(realPath)
	if err != nil {
		return nil, err
	}
	defer closeNoError(file)
	return ioutil.ReadAll(file)
}

func (d *driver) PutContent(subPath string, content []byte) error {
	realPath := d.fullPath(subPath)
	writer, err := d.client.Create(realPath)
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	defer closeNoError(writer)
	return err
}

func (d *driver) Reader(subPath string, offset int64) (io.ReadCloser, error) {
	realPath := d.fullPath(subPath)
	reader, err := d.client.Open(realPath)
	if err != nil {
		return nil, err
	}

	if offset > 0 {
		_, err := reader.Seek(offset, 0)
		if err != nil {
			return nil, err
		}
	}
	return reader, nil
}

func (d *driver) Writer(subPath string, append bool) (io.WriteCloser, error) {
	realPath := d.fullPath(subPath)
	if append {
		return d.client.Append(realPath)
	}
	return d.client.Create(realPath)
}

func (d *driver) Stat(subPath string) (storage.FileInfo, error) {
	realPath := d.fullPath(subPath)
	return d.client.Stat(realPath)
}

func (d *driver) List(subPath string) ([]string, error) {
	realPath := d.fullPath(subPath)
	fileInfos, err := d.client.ReadDir(realPath)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		files = append(files, path.Join(d.rootDirectory, fileInfo.Name()))
	}
	return files, nil
}

func (d *driver) Move(sourcePath, destPath string) error {
	realSrcPath := d.fullPath(sourcePath)
	realDestPath := d.fullPath(destPath)

	err := d.client.Rename(realSrcPath, realDestPath)
	if err != nil {
		if v, ok := err.(*os.PathError); ok {
			if os.IsNotExist(v.Unwrap()) {
				parent := path.Dir(realDestPath)
				err = d.client.MkdirAll(parent, 0755)
				if err != nil {
					return err
				}
				err = d.client.Rename(realSrcPath, realDestPath)
			}
		}
	}
	return err
}

func (d *driver) Delete(path string) error {
	realPath := d.fullPath(path)
	return d.client.RemoveAll(realPath)
}

func (d *driver) Walk(path string, f storage.WalkFn) error {
	return factory.Walk(d, path, f)
}

func closeNoError(closer io.Closer) {
	closer.Close()
}

func (d *driver) fullPath(subPath string) string {
	return path.Join(d.rootDirectory, subPath)
}
