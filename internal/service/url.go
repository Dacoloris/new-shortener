package service

import (
	"context"
	"math/rand"
	"new-shortner/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UrlRepository interface {
	Create(ctx context.Context, book domain.URL) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.URL, error)
	GetAll(ctx context.Context) ([]domain.URL, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Urls struct {
	repo UrlRepository
}

func NewUrls(repo UrlRepository) *Urls {
	return &Urls{
		repo: repo,
	}
}

func (u *Urls) Create(ctx context.Context, url domain.URL) error {
	return u.repo.Create(ctx, url)
}

func (u *Urls) GetByID(ctx context.Context, id uuid.UUID) (domain.URL, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Urls) GetAll(ctx context.Context) ([]domain.URL, error) {
	return u.repo.GetAll(ctx)
}

func (u *Urls) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *Urls) ShortenUrl(_ context.Context) string {
	src := rand.NewSource(time.Now().UnixNano())
	return GenerateURLToken(10, src)
}

// GenerateURLToken generates random base64URL string by given length
// taken from https://stackoverflow.com/a/31832326
func GenerateURLToken(n int, src rand.Source) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
