package file

import (
	"context"
	"encoding/json"
	"io"
	"new-shortner/internal/domain"
	"new-shortner/internal/repository/inmemory"
	"os"

	"go.uber.org/zap"
)

type Storage struct {
	memoryMap *inmemory.URLs
	encoder   *json.Encoder
}

func New(filename string, logger *zap.Logger) (*Storage, error) {
	storage := &Storage{
		memoryMap: inmemory.NewURLs(logger),
	}

	err := storage.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	storage.encoder = json.NewEncoder(file)

	return storage, nil
}

func (s *Storage) Create(ctx context.Context, url domain.URL) error {
	err := s.memoryMap.Create(ctx, url)
	if err != nil {
		return err
	}
	record := struct {
		domain.URL
		UserID string `json:"user_id"`
	}{}
	err = s.encoder.Encode(record)

	return err
}

func (s *Storage) GetOriginalByShort(ctx context.Context, short string) (string, error) {
	return s.memoryMap.GetOriginalByShort(ctx, short)
}

func (s *Storage) GetAllURLsByUserID(ctx context.Context, id string) ([]domain.URL, error) {
	return s.memoryMap.GetAllURLsByUserID(ctx, id)
}

func (s *Storage) LoadFromFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return s.LoadFromBuff(file)
}

func (s *Storage) LoadFromBuff(buf io.Reader) error {
	decoder := json.NewDecoder(buf)

	for {
		var record struct {
			domain.URL
			UserID string `json:"user_id"`
		}
		if err := decoder.Decode(&record); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		url := domain.URL{
			UserID:   record.UserID,
			Original: record.Original,
			Short:    record.Short,
		}
		s.memoryMap.AddRecordToStorage(url)
	}
}
