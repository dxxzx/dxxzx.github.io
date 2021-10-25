package factory

import (
	"io"
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
	reader, err := driver.Reader(path, offset)
	if err != nil {
		defer p.put(driver)
		return nil, err
	}
	return &proxyReader{
		proxy:  p,
		driver: driver,
		reader: reader,
	}, nil
}

func (p *proxy) Writer(path string, append bool) (io.WriteCloser, error) {
	driver := p.get()
	writer, err := driver.Writer(path, append)
	if err != nil {
		defer p.put(driver)
		return nil, err
	}

	return &proxyWriter{
		proxy:  p,
		driver: driver,
		writer: writer,
	}, nil
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

type proxyReader struct {
	proxy  *proxy
	driver storage.Driver
	reader io.ReadCloser
}

func (r *proxyReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *proxyReader) Close() error {
	defer r.proxy.put(r.driver)
	return r.reader.Close()
}

type proxyWriter struct {
	proxy  *proxy
	driver storage.Driver
	writer io.WriteCloser
}

func (w *proxyWriter) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *proxyWriter) Close() error {
	defer w.proxy.put(w.driver)
	return w.writer.Close()
}
