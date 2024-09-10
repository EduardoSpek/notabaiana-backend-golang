package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type CreateDownload interface {
	Create(download *entity.Download) (*entity.Download, error)
}
