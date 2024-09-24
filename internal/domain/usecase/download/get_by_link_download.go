package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type GetByLinkDownloadUsecase struct {
	DownloadRepository port.GetByLinkDownloadRepository
}

func NewGetByLinkDownloadUsecase(repository port.GetByLinkDownloadRepository) *GetByLinkDownloadUsecase {
	return &GetByLinkDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *GetByLinkDownloadUsecase) GetByLink(link string) (*entity.Download, error) {

	download, err := d.DownloadRepository.GetByLink(link)

	if err != nil {
		return &entity.Download{}, err
	}

	return download, nil
}
