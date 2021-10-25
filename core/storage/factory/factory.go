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
				size:    0,
			}
		})
	} else {
		return nil, fmt.Errorf("no such driver factory: %s", name)
	}

	return driver.tryNew(parameters)
}
