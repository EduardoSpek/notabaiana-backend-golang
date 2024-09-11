package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type CreateDownloadRepository interface {
	Create(download *entity.Download) (*entity.Download, error)
}

type UpdateDownloadRepository interface {
	Update(download *entity.Download) (*entity.Download, error)
}
