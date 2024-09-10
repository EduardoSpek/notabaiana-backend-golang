package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type CreateDownloadUsecase struct {
	CreateDownload port.CreateDownload
}

func NewCreateDownloadUsecase(download port.CreateDownload) *CreateDownloadUsecase {
	return &CreateDownloadUsecase{
		CreateDownload: download,
	}
}

func (d *CreateDownloadUsecase) Create(download *entity.Download) (*entity.Download, error) {
	newDownload := entity.NewDownload(*download)
	created, err := d.CreateDownload.Create(newDownload)

	if err != nil {
		return &entity.Download{}, err
	}

	return created, nil
}
