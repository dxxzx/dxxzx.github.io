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
package storage

import (
	"io"
	"os"
)

type Driver interface {
	GetContent(path string) ([]byte, error)
	PutContent(path string, content []byte) error
	Reader(path string, offset int64) (io.ReadCloser, error)
	Writer(path string, flag int, perm os.FileMode) (io.WriteCloser, error)
	Stat(path string) (os.FileInfo, error)
	Readdir(subPath string) ([]os.FileInfo, error)
	Readdirnames(subPath string) ([]string, error)
	Move(sourcePath, destPath string) error
	Delete(path string) error
	Walk(path string, f WalkFn) error
}

type Factory interface {
	Create(params map[string]interface{}) (Driver, error)
}

type WalkFn func(fileInfo os.FileInfo) error
