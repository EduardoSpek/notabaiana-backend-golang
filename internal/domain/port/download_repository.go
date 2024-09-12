package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type DownloadRepository interface {
	Create(download *entity.Download) (*entity.Download, error)
	Update(download *entity.Download) (*entity.Download, error)
	GetByLink(link string) (*entity.Download, error)
	GetBySlug(slug string) (*entity.Download, error)
	FindAll(page, limit int) ([]*entity.Download, error)
	GetTotalVisible() int
	GetTotalSearch(str_search string) int
	GetTotalFindCategory(category string) int
	Search(page int, str_search string) []*entity.Download
	FindCategory(category string, page int) ([]*entity.Download, error)
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

type GetBySlugDownloadRepository interface {
	GetBySlug(slug string) (*entity.Download, error)
}

type FindAllDownloadRepository interface {
	FindAll(page, limit int) ([]*entity.Download, error)
}

type FindCategoryDownloadRepository interface {
	FindCategory(category string, page int) ([]*entity.Download, error)
}

type GetTotalVisibleDownloadRepository interface {
	GetTotalVisible() int
}

type GetTotalFindCategoryDownloadRepository interface {
	GetTotalFindCategory(category string) int
}

type GetTotalSearchDownloadRepository interface {
	GetTotalSearch(str_search string) int
}

type SearchDownloadRepository interface {
	Search(page int, str_search string) []*entity.Download
}
