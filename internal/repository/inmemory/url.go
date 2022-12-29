package inmemory

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type URLs struct {
	logger  *zap.Logger
	storage map[string]string
	mu      sync.RWMutex
}

func NewURLs(lg *zap.Logger) *URLs {
	return &URLs{
		logger:  lg,
		storage: make(map[string]string),
		mu:      sync.RWMutex{},
	}
}

func (u *URLs) Create(ctx context.Context, original, short string) error {
	u.mu.Lock()
	u.storage[short] = original
	u.mu.Unlock()

	return nil
}

func (u *URLs) GetOriginalByShort(ctx context.Context, short string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	url, ok := u.storage[short]
	if ok {
		return url, nil
	}

	return "", errors.New("not found")
}

func (u *URLs) AddRecordToStorage(original, short string) {
	u.mu.Lock()
	u.storage[short] = original
	u.mu.Unlock()
}
