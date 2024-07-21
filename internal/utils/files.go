package utils

import (
	"os"
	"sync"
)

var mutex sync.Mutex

func FileExsists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func RemoveImage(file string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	err := os.Remove(file)
	if err != nil {
		return false
	} else {
		return true
	}
}
