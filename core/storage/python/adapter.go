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
package main

import "C"
import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/dxxzx/magnifier/core/storage"
	"github.com/dxxzx/magnifier/core/storage/factory"
	_ "github.com/dxxzx/magnifier/core/storage/filesystem"
	"gopkg.in/yaml.v2"
)

var driverHolder struct {
	lock    sync.Mutex
	drivers map[C.int]storage.Driver
	idx     C.int
}
var readerHolder struct {
	lock    sync.Mutex
	readers map[C.int]io.ReadCloser
	idx     C.int
}
var writerHolder struct {
	lock    sync.Mutex
	writers map[C.int]io.WriteCloser
	idx     C.int
}

func init() {
	driverHolder.lock.Lock()
	defer driverHolder.lock.Unlock()
	driverHolder.drivers = make(map[C.int]storage.Driver)
	driverHolder.idx = 0

	readerHolder.lock.Lock()
	defer readerHolder.lock.Unlock()
	readerHolder.readers = make(map[C.int]io.ReadCloser)
	readerHolder.idx = 0

	writerHolder.lock.Lock()
	defer writerHolder.lock.Unlock()
	writerHolder.writers = make(map[C.int]io.WriteCloser)
	writerHolder.idx = 0
}

//export CreateStorageDriver
func CreateStorageDriver(name *C.char, parameters *C.char) C.int {
	raw_params := C.GoString(parameters)
	params := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(raw_params), params)
	if err != nil {
		panic(err.Error())
	}
	driver, err := factory.Create(C.GoString(name), params)
	if err != nil {
		panic(err.Error())
	}
	driverHolder.lock.Lock()
	defer driverHolder.lock.Unlock()
	idx := driverHolder.idx
	driverHolder.idx++
	driverHolder.drivers[idx] = driver

	return idx
}

func getDriver(idx C.int) storage.Driver {
	driverHolder.lock.Lock()
	defer driverHolder.lock.Unlock()
	return driverHolder.drivers[idx]
}

//export ReleaseDriver
func ReleaseDriver(idx C.int) {
	driverHolder.lock.Lock()
	defer driverHolder.lock.Unlock()
	delete(driverHolder.drivers, idx)
}

//export GetContent
func GetContent(idx C.int, path *C.char) unsafe.Pointer {
	driver := getDriver(idx)
	data, err := driver.GetContent(C.GoString(path))
	if err != nil {
		panic(err.Error())
	}
	return C.CBytes(data)
}

//export PutContent
func PutContent(idx C.int, path *C.char, content unsafe.Pointer, len C.int) {
	driver := getDriver(idx)
	err := driver.PutContent(C.GoString(path), C.GoBytes(content, len))
	if err != nil {
		panic(err)
	}
}

//export Reader
func Reader(idx C.int, path *C.char, offset C.longlong) C.int {
	driver := getDriver(idx)
	reader, err := driver.Reader(C.GoString(path), int64(offset))
	if err != nil {
		panic(err)
	}
	readerHolder.lock.Lock()
	defer readerHolder.lock.Unlock()
	ridx := readerHolder.idx
	readerHolder.idx++
	readerHolder.readers[ridx] = reader
	return ridx
}

//export Writer
func Writer(idx C.int, path *C.char, append C._Bool) C.int {
	driver := getDriver(idx)
	writer, err := driver.Writer(C.GoString(path), bool(append))
	if err != nil {
		panic(err)
	}
	writerHolder.lock.Lock()
	defer writerHolder.lock.Unlock()
	widx := writerHolder.idx
	writerHolder.idx++
	writerHolder.writers[widx] = writer
	return widx
}

//export Stat
func Stat(idx C.int, path *C.char) *C.char {
	driver := getDriver(idx)
	stat, err := driver.Stat(C.GoString(path))
	if err != nil {
		panic(err)
	}
	resultStr := fmt.Sprintf(
		"%s:%d:%d:%t",
		stat.Name(),
		stat.ModTime().Unix(),
		stat.Size(),
		stat.IsDir(),
	)
	return C.CString(resultStr)
}

//export List
func List(idx C.int, path *C.char) *C.char {
	driver := getDriver(idx)
	result, err := driver.List(C.GoString(path))
	if err != nil {
		panic(err)
	}
	resultStr := strings.Join(result, "\n")
	return C.CString(resultStr)
}

//export Move
func Move(idx C.int, sourcePath *C.char, destPath *C.char) {
	driver := getDriver(idx)
	err := driver.Move(C.GoString(sourcePath), C.GoString(destPath))
	if err != nil {
		panic(err)
	}
}

//export Delete
func Delete(idx C.int, path *C.char) {
	driver := getDriver(idx)
	err := driver.Delete(C.GoString(path))
	if err != nil {
		panic(err)
	}
}

func getReader(idx C.int) io.ReadCloser {
	readerHolder.lock.Lock()
	defer readerHolder.lock.Unlock()
	return readerHolder.readers[idx]
}

//export Read
func Read(idx C.int, buffer unsafe.Pointer, size C.int) C.int {
	reader := getReader(idx)
	var bufHolder []byte
	bufPtr := (*reflect.SliceHeader)(unsafe.Pointer(&bufHolder))
	bufPtr.Data = uintptr(buffer)
	bufPtr.Len = int(size)
	bufPtr.Cap = int(size)
	n, err := reader.Read(bufHolder)
	if err != nil {
		if n == 0 && err != io.EOF {
			panic(err)
		}
	}
	return C.int(n)
}

//export CloseReader
func CloseReader(idx C.int) {
	reader := getReader(idx)
	if reader != nil {
		defer reader.Close()
	}
	readerHolder.lock.Lock()
	defer readerHolder.lock.Unlock()
	delete(readerHolder.readers, idx)
}

func getWriter(idx C.int) io.WriteCloser {
	writerHolder.lock.Lock()
	defer writerHolder.lock.Unlock()
	return writerHolder.writers[idx]
}

//export Write
func Write(idx C.int, buf unsafe.Pointer, size C.int) C.int {
	writer := getWriter(idx)
	n, err := writer.Write(C.GoBytes(buf, size))
	if err != nil {
		panic(err)
	}
	return C.int(n)
}

//export CloseWriter
func CloseWriter(idx C.int) {
	writer := getWriter(idx)
	if writer != nil {
		defer writer.Close()
	}
	writerHolder.lock.Lock()
	defer writerHolder.lock.Unlock()
	delete(writerHolder.writers, idx)
}

func main() {
}
