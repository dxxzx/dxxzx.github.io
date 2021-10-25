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
package hdfs

import (
	"testing"

	"github.com/dxxzx/magnifier/core/storage/factory"
)

const (
	testPath    = "test1.txt"
	testContent = "test content"
)

var parameters = map[string]interface{}{
	"root":      "/user/magnifier",
	"user":      "magnifier",
	"namenodes": []interface{}{"localhost:8020"},
}

func TestGetContent(t *testing.T) {
	driver, err := factory.Create(driverName, parameters)
	if err != nil {
		t.Error(err)
	}
	err = driver.PutContent(testPath, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	data, err := driver.GetContent(testPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != testContent {
		t.Fatal(string(data))
	}
	driver.Delete(testPath)
}

func TestPutContent(t *testing.T) {
	driver, err := factory.Create(driverName, parameters)
	if err != nil {
		t.Error(err)
	}
	err = driver.PutContent(testPath, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	driver.Delete(testPath)
}

func TestMove(t *testing.T) {
	const targetPath = "target/test2.txt"
	driver, err := factory.Create(driverName, parameters)
	if err != nil {
		t.Error(err)
	}

	err = driver.PutContent(testPath, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}

	err = driver.Move(testPath, targetPath)
	if err != nil {
		driver.Delete(testPath)
		t.Fatal(err)
	} else {
		driver.Delete(targetPath)
	}
}
