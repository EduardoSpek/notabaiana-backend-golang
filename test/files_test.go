package test

import (
	"fmt"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

func TestFileExists(t *testing.T) {
	filename := "./images/0aa774a9-a897-4623-afdb-6028aa629ba0.jpg"

	exists := utils.FileExsists(filename)

	fmt.Println(exists)

	if !exists {
		t.Error("O arquivo não existe!")
	}

}