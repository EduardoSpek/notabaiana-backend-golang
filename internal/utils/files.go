package utils

import "os"

func FileExsists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func RemoveImage(file string) bool {
	err := os.Remove(file)
	if err != nil {
		return false
	} else {
		return true
	}
}
