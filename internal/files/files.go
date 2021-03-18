package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var (
	BrevDirectory      = ".brev"
	ActiveProjectsFile = "active_projects.json"
	ProjectsFile       = "projects.json"
	EndpointsFile      = "endpoints.json"
)

func GetBrevDirectory() string {
	return BrevDirectory
}

func GetActiveProjectFile() string {
	return ActiveProjectsFile
}
func GetProjectsFile() string {
	return ProjectsFile
}
func GetEndpointsFile() string {
	return EndpointsFile
}

func GetRootDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func GetActiveProjectsPath() string {
	rootDir := GetRootDir()

	return fmt.Sprintf("%s/%s/%s", rootDir, BrevDirectory, ActiveProjectsFile)
}

func GetLocalBrevDir() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s", cwd, BrevDirectory)
}

func GetEndpointsPath() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s/%s", cwd, BrevDirectory, EndpointsFile)
}
func GetProjectsPath() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s/%s", cwd, BrevDirectory, ProjectsFile)
}

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
