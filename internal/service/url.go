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
	CreateBatch(ctx context.Context, urls []domain.URL) error
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

func (u *URLs) GetAllURLsByUserID(ctx context.Context, userID string) ([]domain.URL, error) {
	if userID == "" {
		return []domain.URL{}, nil
	}
	return u.repo.GetAllURLsByUserID(ctx, userID)
}

func (u *URLs) CreateBatch(
	ctx context.Context,
	req []domain.BatchPostRequest,
	userID string,
) ([]domain.BatchPostResponse, error) {

	Urls := make([]domain.URL, 0, len(req))
	src := rand.NewSource(time.Now().UnixNano())

	for _, elem := range req {
		var url domain.URL
		url.UserID = userID
		url.Original = elem.Original
		url.Short = GenerateURLToken(10, src)
		Urls = append(Urls, url)
	}

	err := u.repo.CreateBatch(ctx, Urls)
	if err != nil {
		return nil, err
	}

	res := make([]domain.BatchPostResponse, 0, len(req))

	for i := 0; i < len(req); i++ {
		var elem domain.BatchPostResponse
		elem.CorrelationID = req[i].CorrelationID
		elem.Short = Urls[i].Short
		res = append(res, elem)
	}

	return res, nil
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
