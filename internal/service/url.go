package service

import (
	"context"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type URLRepository interface {
	Create(ctx context.Context, original, short string) error
	GetOriginalByShort(ctx context.Context, short string) (string, error)
}

type URLs struct {
	repo URLRepository
}

func NewURLs(repo URLRepository) *URLs {
	return &URLs{
		repo: repo,
	}
}

func (u *URLs) Create(ctx context.Context, original string) (string, error) {
	_, err := url.ParseRequestURI(original)
	if err != nil {
		return "", err
	}

	src := rand.NewSource(time.Now().UnixNano())
	short := GenerateURLToken(10, src)

	err = u.repo.Create(ctx, original, short)
	if err != nil {
		return "", err
	}

	return short, nil
}

func (u *URLs) GetOriginalByShort(ctx context.Context, short string) (string, error) {
	return u.repo.GetOriginalByShort(ctx, short)
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
