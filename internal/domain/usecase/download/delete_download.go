package usecase

import (
	"fmt"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type DeleteDownloadUsecase struct {
	DownloadRepository port.DeleteDownloadRepository
}

func NewDeleteDownloadUsecase(repository port.DeleteDownloadRepository) *DeleteDownloadUsecase {
	return &DeleteDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *DeleteDownloadUsecase) Delete(id string) error {

	download, err := d.DownloadRepository.GetByID(id)

	if err != nil {
		return err
	}

	removed := utils.RemoveImage("images/downloads/" + download.Image)

	if !removed {
		fmt.Println("Delete Download: não foi possível deletar a imagem")
	}

	err = d.DownloadRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
