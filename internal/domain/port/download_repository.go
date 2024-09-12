package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type DownloadRepository interface {
	Create(download *entity.Download) (*entity.Download, error)
	Update(download *entity.Download) (*entity.Download, error)
	GetByLink(link string) (*entity.Download, error)
}

type CreateDownloadRepository interface {
	Create(download *entity.Download) (*entity.Download, error)
}

type UpdateDownloadRepository interface {
	Update(download *entity.Download) (*entity.Download, error)
}

type GetByLinkDownloadRepository interface {
	GetByLink(download string) (*entity.Download, error)
}
