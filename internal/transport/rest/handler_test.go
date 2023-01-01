package rest

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"new-shortner/internal/config"
	"new-shortner/internal/domain"
	"new-shortner/internal/transport/rest/cookie"
	mock_rest "new-shortner/internal/transport/rest/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	type mockBehavior func(m *mock_rest.MockURLs, ctx context.Context, short string)

	testTable := []struct {
		name                string
		method              string
		request             string
		short               string
		mockBehavior        mockBehavior
		exceptedContentType string
		exceptedStatusCode  int
		exceptedLocation    string
	}{
		{
			name:    "OK",
			method:  http.MethodGet,
			request: "/short",
			short:   "short",
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, short string) {
				m.EXPECT().GetOriginalByShort(ctx, short).Return("https://google.com", nil)
			},
			exceptedContentType: "text/plain",
			exceptedStatusCode:  http.StatusTemporaryRedirect,
			exceptedLocation:    "https://google.com",
		},
		{
			name:    "invalid request uri",
			method:  http.MethodGet,
			request: "/s",
			short:   "s",
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, short string) {
				m.EXPECT().GetOriginalByShort(ctx, short).Return("", errors.New("not found"))
			},
			exceptedContentType: "text/plain",
			exceptedStatusCode:  http.StatusBadRequest,
			exceptedLocation:    "",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			repo := mock_rest.NewMockURLs(ctrl)
			tt.mockBehavior(repo, ctx, tt.short)

			cfg, err := config.New()
			assert.NoError(t, err)
			h := NewHandler(repo, cfg)

			router := gin.New()
			router.Use(SetPlainTextHeader(), cookie.CheckCookie)
			router.GET("/:id", h.Redirect)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.request, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.exceptedContentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.exceptedStatusCode, w.Code)
			assert.Equal(t, tt.exceptedLocation, w.Header().Get("Location"))
		})
	}
}

func TestURLShortening(t *testing.T) {
	type mockBehavior func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL)

	testTable := []struct {
		name                string
		method              string
		url                 domain.URL
		request             string
		requestBody         string
		mockBehavior        mockBehavior
		exceptedContentType string
		exceptedStatusCode  int
		exceptedBody        string
	}{
		{
			name:        "OK",
			method:      http.MethodPost,
			url:         domain.URL{Original: "http://google.com"},
			request:     "/",
			requestBody: "http://google.com",
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL) {
				m.EXPECT().Create(ctx, url).Return("short", nil)
			},
			exceptedContentType: "text/plain",
			exceptedStatusCode:  http.StatusCreated,
			exceptedBody:        "http://localhost:8080/short",
		},
		{
			name:        "invalid body url",
			method:      http.MethodPost,
			url:         domain.URL{Original: "http/google.com"},
			request:     "/",
			requestBody: "http/google.com",
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL) {
				m.EXPECT().Create(ctx, url).Return("", errors.New("parse uri fail"))
			},
			exceptedContentType: "text/plain",
			exceptedStatusCode:  http.StatusBadRequest,
			exceptedBody:        "parse uri fail",
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			url:         domain.URL{Original: ""},
			request:     "/",
			requestBody: "",
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL) {
				m.EXPECT().Create(ctx, url).Return("", errors.New("parse uri fail"))
			},
			exceptedContentType: "text/plain",
			exceptedStatusCode:  http.StatusBadRequest,
			exceptedBody:        "parse uri fail",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			repo := mock_rest.NewMockURLs(ctrl)
			tt.mockBehavior(repo, ctx, tt.url)

			cfg, err := config.New()
			assert.NoError(t, err)
			h := NewHandler(repo, cfg)

			router := gin.New()
			router.Use(SetPlainTextHeader(), cookie.CheckCookie)
			router.POST("/", h.URLShortening)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.request, bytes.NewBufferString(tt.requestBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.exceptedContentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.exceptedStatusCode, w.Code)
			assert.Equal(t, tt.exceptedBody, w.Body.String())
		})
	}
}

func TestAPIShorten(t *testing.T) {
	type mockBehavior func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL)

	testTable := []struct {
		name                string
		method              string
		url                 domain.URL
		request             string
		requestBody         string
		mockBehavior        mockBehavior
		exceptedContentType string
		exceptedStatusCode  int
		exceptedBody        string
	}{
		{
			name:        "OK",
			method:      http.MethodPost,
			url:         domain.URL{Original: "http://google.com"},
			request:     "/api/shorten",
			requestBody: `{"url":"http://google.com"}`,
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL) {
				m.EXPECT().Create(ctx, url).Return("short", nil)
			},
			exceptedContentType: "application/json",
			exceptedStatusCode:  http.StatusCreated,
			exceptedBody:        `{"result":"http://localhost:8080/short"}`,
		},
		{
			name:        "invalid body url",
			method:      http.MethodPost,
			url:         domain.URL{Original: "ht//google"},
			request:     "/api/shorten",
			requestBody: `{"url":"ht//google"}`,
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL) {
				m.EXPECT().Create(ctx, url).Return("", errors.New("invalid url"))
			},
			exceptedContentType: "application/json",
			exceptedStatusCode:  http.StatusBadRequest,
			exceptedBody:        `{"result":"invalid url"}`,
		},
		{
			name:        "invalid json",
			method:      http.MethodPost,
			url:         domain.URL{Original: "http://google.com"},
			request:     "/api/shorten",
			requestBody: `{"u":"http://google.com"}`,
			mockBehavior: func(m *mock_rest.MockURLs, ctx context.Context, url domain.URL) {
			},
			exceptedContentType: "application/json",
			exceptedStatusCode:  http.StatusBadRequest,
			exceptedBody:        `{"result":"invalid json"}`,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			repo := mock_rest.NewMockURLs(ctrl)
			tt.mockBehavior(repo, ctx, tt.url)

			cfg, err := config.New()
			assert.NoError(t, err)
			h := NewHandler(repo, cfg)

			router := gin.New()
			router.Use(SetJSONHeader())
			router.POST("/api/shorten", h.APIShorten)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.request, bytes.NewBufferString(tt.requestBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.exceptedContentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.exceptedStatusCode, w.Code)
			assert.Equal(t, tt.exceptedBody, w.Body.String())
		})
	}
}
