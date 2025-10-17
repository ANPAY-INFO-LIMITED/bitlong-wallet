package base

import (
	"fmt"
	"os"
	"sync"
)

var (
	mu       sync.Mutex
	filePath string
)

func SetFilePath(path string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist at path: %s", path)
	}
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

func GetFilePath() string {
	mu.Lock()
	defer mu.Unlock()
	return filePath
}
