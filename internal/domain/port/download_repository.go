package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type DownloadRepository interface {
	Create(download *entity.Download) (*entity.Download, error)
	Update(download *entity.Download) (*entity.Download, error)
	GetByID(id string) (*entity.Download, error)
	GetByLink(link string) (*entity.Download, error)
	GetBySlug(slug string) (*entity.Download, error)
	AdminGetBySlug(slug string) (*entity.Download, error)
	FindAll(page, limit int) ([]*entity.Download, error)
	FindAllTopViews(page, limit int) ([]*entity.Download, error)
	AdminFindAll(page, limit int) ([]*entity.Download, error)
	GetTotal() (int, error)
	GetTotalVisible() (int, error)
	GetTotalSearch(str_search string) (int, error)
	GetTotalFindCategory(category string) (int, error)
	Search(page int, str_search string) ([]*entity.Download, error)
	FindCategory(category string, page int) ([]*entity.Download, error)
	Delete(id string) error
	DeleteAll(downloads []*entity.Download) error
	Clean() ([]*entity.Download, error)
}

type CopierDownload interface {
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

type GetByIDDownloadRepository interface {
	GetByID(id string) (*entity.Download, error)
}

type GetByLinkDownloadRepository interface {
	GetByLink(download string) (*entity.Download, error)
}

type GetBySlugDownloadRepository interface {
	GetBySlug(slug string) (*entity.Download, error)
}

type AdminGetBySlugDownloadRepository interface {
	AdminGetBySlug(slug string) (*entity.Download, error)
}

type FindAllDownloadRepository interface {
	FindAll(page, limit int) ([]*entity.Download, error)
	GetTotalVisible() (int, error)
}

type AdminFindAllDownloadRepository interface {
	AdminFindAll(page, limit int) ([]*entity.Download, error)
	GetTotal() (int, error)
}

type FindAllTopViewsDownloadRepository interface {
	FindAllTopViews(page, limit int) ([]*entity.Download, error)
}

type FindCategoryDownloadRepository interface {
	FindCategory(category string, page int) ([]*entity.Download, error)
	GetTotalFindCategory(category string) (int, error)
}

type GetTotalDownloadRepository interface {
	GetTotal() (int, error)
}

type GetTotalVisibleDownloadRepository interface {
	GetTotalVisible() (int, error)
}

type GetTotalFindCategoryDownloadRepository interface {
	GetTotalFindCategory(category string) (int, error)
}

type GetTotalSearchDownloadRepository interface {
	GetTotalSearch(str_search string) (int, error)
}

type SearchDownloadRepository interface {
	Search(page int, str_search string) ([]*entity.Download, error)
	GetTotalSearch(str_search string) (int, error)
}

type DeleteDownloadRepository interface {
	Delete(id string) error
	GetByID(id string) (*entity.Download, error)
}

type DeleteAllDownloadRepository interface {
	DeleteAll(downloads []*entity.Download) error
	GetByID(id string) (*entity.Download, error)
}

type CleanDownloadRepository interface {
	Clean() ([]*entity.Download, error)
	DeleteAll(downloads []*entity.Download) error
}
