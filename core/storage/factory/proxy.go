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
package factory

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/dxxzx/magnifier/core/storage"
)

type proxy struct {
	drivers []storage.Driver
	limit   int32
	factory storage.Factory
	mutex   sync.Mutex
	size    int32
}

func (p *proxy) get() storage.Driver {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	length := len(p.drivers)
	driver := p.drivers[length-1]
	p.drivers = p.drivers[:length-1]
	p.size--
	return driver
}

func (p *proxy) put(driver storage.Driver) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.drivers = append(p.drivers, driver)
	p.size++
}

func (p *proxy) tryNew(parameters map[string]interface{}) (storage.Driver, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.size < p.limit || len(p.drivers) == 0 {
		driver, err := p.factory.Create(parameters)
		if err != nil {
			return nil, err
		}
		p.drivers = append(p.drivers, driver)
		p.size++
	}
	return p, nil
}

func (p *proxy) GetContent(path string) ([]byte, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.GetContent(path)
}
func (p *proxy) PutContent(path string, content []byte) error {
	driver := p.get()
	defer p.put(driver)
	return driver.PutContent(path, content)
}

func (p *proxy) Reader(path string, offset int64) (io.ReadCloser, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Reader(path, offset)
}

func (p *proxy) Writer(path string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Writer(path, flag, perm)
}

func (p *proxy) Stat(path string) (os.FileInfo, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Stat(path)
}

func (p *proxy) Readdir(path string) ([]os.FileInfo, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Readdir(path)
}

func (p *proxy) Readdirnames(path string) ([]string, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Readdirnames(path)
}

func (p *proxy) Move(sourcePath, destPath string) error {
	driver := p.get()
	defer p.put(driver)
	return driver.Move(sourcePath, destPath)
}

func (p *proxy) Delete(path string) error {
	driver := p.get()
	defer p.put(driver)
	return driver.Delete(path)
}

func (p *proxy) Walk(path string, fn storage.WalkFn) error {
	driver := p.get()
	defer p.put(driver)
	return Walk(driver, path, fn)
}

func Walk(driver storage.Driver, subPath string, fn storage.WalkFn) error {
	children, err := driver.Readdir(subPath)
	if err != nil {
		return err
	}

	for _, child := range children {
		err = fn(child)
		if err != nil {
			return err
		}
		if child.IsDir() {
			if err := Walk(driver, path.Join(subPath, child.Name()), fn); err != nil {
				return err
			}
		}
	}
	return nil
}
