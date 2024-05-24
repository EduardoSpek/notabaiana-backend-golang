package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type HitsRepository interface {
	Save(hit entity.Hits) error
	Update(hit entity.Hits) error
	Get(ip string, session string) (entity.Hits, error)
	TopHits() ([]entity.Hits, error)
}