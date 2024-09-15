package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type AdminGetBySlugDownloadUsecase struct {
	DownloadRepository port.AdminGetBySlugDownloadRepository
}

func NewAdminGetBySlugDownloadUsecase(repository port.AdminGetBySlugDownloadRepository) *AdminGetBySlugDownloadUsecase {
	return &AdminGetBySlugDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *AdminGetBySlugDownloadUsecase) AdminGetBySlug(slug string) (*entity.Download, error) {

	download, err := d.DownloadRepository.AdminGetBySlug(slug)

	if err != nil {
		return &entity.Download{}, err
	}

	return download, nil
}
