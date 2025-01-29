package usecase

import (
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type CleanOldDownloadUsecase struct {
	DownloadRepository port.CleanOldDownloadRepository
}

func NewCleanOldDownloadUsecase(repository port.CleanOldDownloadRepository) *CleanOldDownloadUsecase {
	return &CleanOldDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *CleanOldDownloadUsecase) StartCleanOldDownloads(minutes time.Duration) {

	go d.CleanOld()

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go d.CleanOld()
	}
}

func (d *CleanOldDownloadUsecase) CleanOld() error {

	downloads, err := d.DownloadRepository.CleanOld()

	if err != nil {
		return err
	}

	for _, n := range downloads {

		if n.Image != "" {
			image := "./downloads/" + n.Image
			utils.RemoveImage(image)
		}
	}

	err = d.DownloadRepository.DeleteAll(downloads)

	if err != nil {
		return err
	}

	return nil

}
