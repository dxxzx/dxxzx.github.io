package filesystem

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dxxzx/magnifier/core/storage"
)

const (
	testPath1   = "testPath1.txt"
	testContent = "test Content"
	testPath2   = "testPath2.txt"
)

func getDriver(t *testing.T) *driver {
	params := map[string]interface{}{
		"root": "/tmp/magnifier",
	}
	driver, err := FromParameters(params)
	if err != nil {
		t.Fatalf("get driver failed: %s", err.Error())
	}
	return driver
}

func TestGetContent(t *testing.T) {
	driver := getDriver(t)
	err := driver.PutContent(testPath1, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath1)

	content, err := driver.GetContent(testPath1)
	if string(content) != testContent {
		t.Fatalf("content doesn't match: %s", err.Error())
	}
}

func TestStat(t *testing.T) {
	driver := getDriver(t)
	err := driver.PutContent(testPath1, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath1)

	fileInfo, err := driver.Stat(testPath1)
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.IsDir() {
		t.Fatal("target path is not dir")
	}
	if fileInfo.Name() != testPath1 {
		t.Fatalf("error filename")
	}
	if fileInfo.Size() != 12 {
		t.Fatalf("error file size")
	}
}

func TestList(t *testing.T) {
	driver := getDriver(t)
	err := driver.PutContent(testPath1, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath1)

	err = driver.PutContent(testPath2, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath2)

	results, err := driver.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatal("error file count")
	}
	if !reflect.DeepEqual(results, []string{testPath1, testPath2}) {
		t.Fatalf("error result: %v", results)
	}
}
func TestMove(t *testing.T) {
	driver := getDriver(t)
	err := driver.PutContent(testPath1, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath1)

	err = driver.Move(testPath1, testPath2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = driver.Stat(testPath1)
	if err == nil {
		t.Fatalf("origin file still exists")
	}

	content, err := driver.GetContent(testPath2)
	if string(content) != testContent {
		t.Fatalf("content doesn't match: %s", err.Error())
	}
}

func TestWalk(t *testing.T) {
	driver := getDriver(t)
	err := driver.PutContent(testPath1, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath1)

	err = driver.PutContent(testPath2, []byte(testContent))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Delete(testPath2)

	driver.Walk("", func(fileInfo storage.FileInfo) error {
		fmt.Printf("fileinfo: %#v\n", fileInfo)
		return nil
	})
}
