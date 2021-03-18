package files

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Read data into the given struct
//
// Usage:
//   var foo myStruct
//   files.ReadJSON("tmp/a.json", &foo)
func ReadJSON(filepath string, v interface{}) error {
	f, err := os.Open(filepath)
	defer f.Close()

	if err != nil {
		return err
	}

	dataBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(dataBytes, v)
}

// Replace data in the target file with data from the given struct
//
// Usage (unstructured):
//   OverwriteJSON("tmp/a/b/c.json", map[string]string{
// 	    "hi": "there",
//   })
//
//
// Usge (struct):
//   var foo myStruct
//   OverwriteJSON("tmp/a/b/c.json", foo)
func OverwriteJSON(filepath string, v interface{}) error {
	f, err := touchFile(filepath)
	if err != nil {
		return nil
	}
	defer f.Close()

	// clear
	err = f.Truncate(0)

	// write
	dataBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, dataBytes, os.ModePerm)

	return err
}

// Create file (and full path) if it does not already exit
func touchFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func Does_file_exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// file does not exist
		return false
	} else {
		// file exists
		return true
	}
} 