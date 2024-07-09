package base

import (
	"fmt"
	"os"
	"sync"
)

var (
	// Use mutexes to ensure thread-safe access to paths
	mu sync.Mutex
	// The path to the storage file
	filePath string
)

// SetFilePath
// Set the file path and perform some validation at setup time
func SetFilePath(path string) error {
	mu.Lock()
	defer mu.Unlock()
	//fmt.Printf("path:%v\n", path)
	// Here you can add the path validation logic
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist at path: %s", path)
	}
	// Let's say it's just reading the file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file")
		}
	}(file)
	filePath = path
	return nil
}

// GetFilePath
// Get the path to the stored file
func GetFilePath() string {
	mu.Lock()
	defer mu.Unlock()
	return filePath
}
