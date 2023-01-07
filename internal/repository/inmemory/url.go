package inmemory

import (
	"context"
	"errors"
	"new-shortner/internal/domain"
	"sync"

	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("not found")
)

type URLs struct {
	logger  *zap.Logger
	storage map[string]domain.URL
	mu      sync.RWMutex
}

func NewURLs(lg *zap.Logger) *URLs {
	return &URLs{
		logger:  lg,
		storage: make(map[string]domain.URL),
		mu:      sync.RWMutex{},
	}
}

func (u *URLs) Create(_ context.Context, url domain.URL) (string, error) {
	u.AddRecordToStorage(url)
	return "", nil
}

func (u *URLs) GetOriginalByShort(_ context.Context, short string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	url, ok := u.storage[short]
	if ok {
		return url.Original, nil
	}

	return "", ErrNotFound
}

func (u *URLs) AddRecordToStorage(url domain.URL) {
	u.mu.Lock()
	u.storage[url.Short] = url
	u.mu.Unlock()
}

func (u *URLs) GetAllURLsByUserID(_ context.Context, userID string) ([]domain.URL, error) {
	res := make([]domain.URL, 0)
	u.mu.RLock()
	for _, url := range u.storage {
		if url.UserID == userID {
			res = append(res, url)
		}
	}
	u.mu.RUnlock()

	return res, nil
}

func (u *URLs) CreateBatch(_ context.Context, urls []domain.URL) error {
	for _, url := range urls {
		u.AddRecordToStorage(url)
	}

	return nil
}
