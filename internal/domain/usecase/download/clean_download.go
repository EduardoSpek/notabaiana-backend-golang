package usecase

import (
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type CleanDownloadUsecase struct {
	DownloadRepository port.CleanDownloadRepository
}

func NewCleanDownloadUsecase(repository port.CleanDownloadRepository) *CleanDownloadUsecase {
	return &CleanDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *CleanDownloadUsecase) StartCleanDownloads(minutes time.Duration) {

	go d.Clean()

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go d.Clean()
	}
}

func (d *CleanDownloadUsecase) Clean() error {

	downloads, err := d.DownloadRepository.Clean()

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
