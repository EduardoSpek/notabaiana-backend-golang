package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type UpdateDownloadUsecase struct {
	DownloadRepository port.UpdateDownloadRepository
}

func NewUpdateDownloadUsecase(repository port.UpdateDownloadRepository) *UpdateDownloadUsecase {
	return &UpdateDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *UpdateDownloadUsecase) Update(download *entity.Download) (*entity.Download, error) {
	updateDownload := entity.UpdateDownload(*download)

	updated, err := d.DownloadRepository.Update(updateDownload)

	if err != nil {
		return &entity.Download{}, err
	}

	return updated, nil
}
