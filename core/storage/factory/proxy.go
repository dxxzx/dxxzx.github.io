package factory

import (
	"io"

	"github.com/dxxzx/magnifier/core/storage"
)

type Proxy struct {
	drivers []storage.Driver
	limit   int32
	factory storage.Factory
	lock    chan struct{}
	size    int32
}

func (p *Proxy) get() storage.Driver {
	p.lock <- struct{}{}
	length := len(p.drivers)
	driver := p.drivers[length-1]
	p.drivers = p.drivers[:length-1]
	p.size--
	<-p.lock
	return driver
}

func (p *Proxy) put(driver storage.Driver) {
	p.lock <- struct{}{}
	p.drivers = append(p.drivers, driver)
	p.size++
	<-p.lock
}

func (p *Proxy) tryNew(parameters map[string]interface{}) (storage.Driver, error) {
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

func (p *Proxy) GetContent(path string) ([]byte, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.GetContent(path)
}
func (p *Proxy) PutContent(path string, content []byte) error {
	driver := p.get()
	defer p.put(driver)
	return driver.PutContent(path, content)
}
func (p *Proxy) Reader(path string, offset int64) (io.ReadCloser, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Reader(path, offset)
}
func (p *Proxy) Writer(path string, append bool) (io.WriteCloser, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Writer(path, append)
}
func (p *Proxy) Stat(path string) (storage.FileInfo, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.Stat(path)
}
func (p *Proxy) List(path string) ([]string, error) {
	driver := p.get()
	defer p.put(driver)
	return driver.List(path)
}
func (p *Proxy) Move(sourcePath, destPath string) error {
	driver := p.get()
	defer p.put(driver)
	return driver.Move(sourcePath, destPath)
}
func (p *Proxy) Delete(path string) error {
	driver := p.get()
	defer p.put(driver)
	return driver.Delete(path)
}
func (p *Proxy) Walk(path string, fn storage.WalkFn) error {
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
