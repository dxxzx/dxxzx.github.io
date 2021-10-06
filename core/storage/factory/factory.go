package factory

import (
	"fmt"
	"sync"

	"github.com/dxxzx/magnifier/core/storage"
)

var (
	factories = make(map[string]storage.Factory)
	driver    *proxy
	once      sync.Once
)

func Register(name string, factory storage.Factory) {
	if factory == nil {
		panic("must not provide nil storage factory")
	}
	_, registered := factories[name]
	if registered {
		panic(fmt.Sprintf("storage factory named %s already registered", name))
	}

	factories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (storage.Driver, error) {
	factory, ok := factories[name]
	if ok {
		once.Do(func() {
			driver = &proxy{
				drivers: make([]storage.Driver, 8),
				limit:   32,
				factory: factory,
				lock:    make(chan struct{}, 1),
				size:    0,
			}
		})
	} else {
		return nil, fmt.Errorf("no such driver factory: %s", name)
	}

	return driver.tryNew(parameters)
}
