package service

import (
	"context"
	"errors"
	"math/rand"
	urls "net/url"
	"new-shortner/internal/domain"
	"strings"
	"time"
)

var (
	ErrParseURI = errors.New("parse uri fail")
)

type URLRepository interface {
	Create(ctx context.Context, url domain.URL) error
	GetOriginalByShort(ctx context.Context, short string) (string, error)
	GetAllURLsByUserID(ctx context.Context, id string) ([]domain.URL, error)
}

type URLs struct {
	repo URLRepository
}

func NewURLs(repo URLRepository) *URLs {
	return &URLs{
		repo: repo,
	}
}

func (u *URLs) Create(ctx context.Context, url domain.URL) (string, error) {
	_, err := urls.ParseRequestURI(url.Original)
	if err != nil {
		return "", ErrParseURI
	}

	src := rand.NewSource(time.Now().UnixNano())
	url.Short = GenerateURLToken(10, src)

	err = u.repo.Create(ctx, url)
	if err != nil {
		return "", err
	}

	return url.Short, nil
}

func (u *URLs) GetOriginalByShort(ctx context.Context, short string) (string, error) {
	return u.repo.GetOriginalByShort(ctx, short)
}

func (u *URLs) GetAllURLsByUserID(ctx context.Context, id string) ([]domain.URL, error) {
	if id == "" {
		return []domain.URL{}, nil
	}
	return u.repo.GetAllURLsByUserID(ctx, id)
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
