package test

import (
	"testing"

	"github.com/eduardospek/bn-api/internal/utils"
)

func TestFileExists(t *testing.T) {
	filename := "/images/0bff42bf-8a1b-4457-a208-4107b0dda8d9.jpg"

	exists := utils.FileExsists(filename)

	if !exists {
		t.Error("O arquivo não existe!")
	}

}