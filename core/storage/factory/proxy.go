package factory

import (
	"io"

	"github.com/dxxzx/magnifier/core/storage"
)

type proxy struct {
	drivers []storage.Driver
	limit   int32
	factory storage.Factory
	lock    chan struct{}
	size    int32
}

func (p *proxy) get() storage.Driver {
	p.lock <- struct{}{}
	length := len(p.drivers)
	driver := p.drivers[length-1]
	p.drivers = p.drivers[:length-1]
	p.size--
	<-p.lock
	return driver
}

func (p *proxy) put(driver storage.Driver) {
	p.lock <- struct{}{}
	p.drivers = append(p.drivers, driver)
	p.size++
	<-p.lock
}

func (p *proxy) tryNew(parameters map[string]interface{}) (storage.Driver, error) {
	p.lock <- struct{}{}
	if p.size < p.limit || len(p.drivers) == 0 {
		driver, err := p.factory.Create(parameters)
		if err != nil {
			<-p.lock
			return nil, err
		}
		p.drivers = append(p.drivers, driver)
		p.size++
	}
	<-p.lock
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
func (p *proxy) Writer(path string, append bool) (io.WriteCloser, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Writer(path, append)
}
func (p *proxy) Stat(path string) (storage.FileInfo, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Stat(path)
}
func (p *proxy) List(path string) ([]string, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.List(path)
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

func Walk(driver storage.Driver, path string, fn storage.WalkFn) error {
	children, err := driver.List(path)
	if err != nil {
		return err
	}

	for _, child := range children {
		fileInfo, err := driver.Stat(child)
		if err != nil {
			return nil
		}

		err = fn(fileInfo)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if err := Walk(driver, child, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
