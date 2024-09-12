package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type DownloadRepository interface {
	Create(download *entity.Download) (*entity.Download, error)
	Update(download *entity.Download) (*entity.Download, error)
	GetByLink(link string) (*entity.Download, error)
	FindAll(page, limit int) ([]*entity.Download, error)
	GetTotalVisible() int
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

type FindAllDownloadRepository interface {
	FindAll(page, limit int) ([]*entity.Download, error)
}

type GetTotalVisibleDownloadRepository interface {
	GetTotalVisible() int
}
