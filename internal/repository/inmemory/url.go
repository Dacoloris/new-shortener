package inmemory

import (
	"context"
	"errors"
	"new-shortner/internal/domain"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Urls struct {
	logger  *zap.Logger
	storage map[uuid.UUID]domain.URL
	mu      sync.RWMutex
}

func NewUrls(lg *zap.Logger) *Urls {
	return &Urls{
		logger:  lg,
		storage: make(map[uuid.UUID]domain.URL),
		mu:      sync.RWMutex{},
	}
}

func (u *Urls) Create(ctx context.Context, url domain.URL) error {
	u.mu.Lock()
	u.storage[url.ID] = url
	u.mu.Unlock()

	return nil
}

func (u *Urls) GetByID(ctx context.Context, id uuid.UUID) (domain.URL, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	url, ok := u.storage[id]
	if ok {
		return url, nil
	}

	return domain.URL{}, errors.New("not found")
}

func (u *Urls) GetAll(ctx context.Context) ([]domain.URL, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	res := make([]domain.URL, 0, len(u.storage))
	for _, v := range u.storage {
		res = append(res, v)
	}

	return res, nil
}

func (u *Urls) Delete(ctx context.Context, id uuid.UUID) error {
	delete(u.storage, id)

	return nil
}
