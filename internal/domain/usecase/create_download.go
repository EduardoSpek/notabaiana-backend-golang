package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type CreateDownloadUsecase struct {
	DownloadRepository port.CreateDownloadRepository
}

func NewCreateDownloadUsecase(repository port.CreateDownloadRepository) *CreateDownloadUsecase {
	return &CreateDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *CreateDownloadUsecase) Create(download *entity.Download) (*entity.Download, error) {
	newDownload := entity.NewDownload(*download)

	newDownload.Image = newDownload.ID + ".jpg"

	created, err := d.DownloadRepository.Create(newDownload)

	if err != nil {
		return &entity.Download{}, err
	}

	return created, nil
}
