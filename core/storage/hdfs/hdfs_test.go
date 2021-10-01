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
