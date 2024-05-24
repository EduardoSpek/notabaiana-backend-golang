package memorydb

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

type HitsMemoryRepository struct {
	HitsDB map[string]entity.Hits
}

func NewHitsMemoryRepository() *HitsMemoryRepository {
	return &HitsMemoryRepository{ HitsDB: make(map[string]entity.Hits) }
}

func (h *HitsMemoryRepository) Save(hit entity.Hits) error {
	return nil
}

func (h *HitsMemoryRepository) Update(hit entity.Hits) error {
	return nil
}

func (h *HitsMemoryRepository) Get(ip string, session string) (entity.Hits, error) {
	return entity.Hits{}, nil
}

func (h *HitsMemoryRepository) TopHits() ([]entity.Hits, error) {
	return []entity.Hits{}, nil
}
