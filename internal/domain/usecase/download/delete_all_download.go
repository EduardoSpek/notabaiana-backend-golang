package usecase

import (
	"fmt"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type DeleteAllDownloadUsecase struct {
	DownloadRepository port.DownloadRepository
}

func NewDeleteAllDownloadUsecase(repository port.DownloadRepository) *DeleteAllDownloadUsecase {
	return &DeleteAllDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *DeleteAllDownloadUsecase) DeleteAll(downloads []*entity.Download) error {

	for _, download := range downloads {

		download, err := d.DownloadRepository.GetByID(download.ID)

		if err != nil {
			return err
		}

		removed := utils.RemoveImage("images/downloads/" + download.Image)

		if !removed {
			fmt.Println("DeleteAll Download: não foi possível deletar a imagem")
		}

	}

	err := d.DownloadRepository.DeleteAll(downloads)

	if err != nil {
		return err
	}

	return nil
}
