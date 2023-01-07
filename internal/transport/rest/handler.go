package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"new-shortner/internal/config"
	"new-shortner/internal/domain"
	"new-shortner/internal/transport/rest/cookie"
	"new-shortner/internal/transport/rest/ping"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSON = errors.New("invalid json")
	ErrInvalidURL  = errors.New("invalid url")
)

type URLs interface {
	Create(ctx context.Context, url domain.URL) (string, error)
	GetOriginalByShort(ctx context.Context, short string) (string, error)
	GetAllURLsByUserID(ctx context.Context, UserID string) ([]domain.URL, error)
	CreateBatch(ctx context.Context, req []domain.BatchPostRequest, userID, baseURL string) ([]domain.BatchPostResponse, error)
}

type Handler struct {
	URLsService URLs
	cfg         config.Config
}

func NewHandler(urls URLs, cfg config.Config) *Handler {
	return &Handler{
		URLsService: urls,
		cfg:         cfg,
	}
}

func (h *Handler) Redirect(c *gin.Context) {
	short := c.Param("id")
	original, err := h.URLsService.GetOriginalByShort(c.Request.Context(), short)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Location", original)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) URLShortening(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	id, err := cookie.ReadEncrypted(c.Request, "id")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := domain.URL{
		UserID:   id,
		Original: string(b),
	}

	short, err := h.URLsService.Create(c.Request.Context(), url)
	if err != nil {
		var uc *domain.UniqueConstraintError
		if errors.As(err, &uc) {
			c.String(http.StatusConflict, fmt.Sprintf("%s/%s", h.cfg.BaseURL, short))
			return
		}
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusCreated, fmt.Sprintf("%s/%s", h.cfg.BaseURL, short))
}

func (h *Handler) APIShorten(c *gin.Context) {
	j := struct {
		URL string `json:"url"`
	}{}
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
		return
	}
	err = json.Unmarshal(b, &j)
	if err != nil || j.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": ErrInvalidJSON.Error()})
		return
	}

	id, err := cookie.ReadEncrypted(c.Request, "id")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	url := domain.URL{
		UserID:   id,
		Original: j.URL,
	}

	short, err := h.URLsService.Create(c.Request.Context(), url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": ErrInvalidURL.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": fmt.Sprintf("%s/%s", h.cfg.BaseURL, short)})
}

func (h *Handler) GetAllURLsForUser(c *gin.Context) {
	id, err := cookie.ReadEncrypted(c.Request, "id")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	urls, err := h.URLsService.GetAllURLsByUserID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(urls) != 0 {
		res := make([]domain.URL, len(urls))
		for i := 0; i < len(urls); i++ {
			res[i] = urls[i]
			res[i].Short = fmt.Sprintf("%s/%s", h.cfg.BaseURL, urls[i].Short)
		}
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusNoContent, urls)
	}
}

func (h *Handler) APIBatch(c *gin.Context) {
	var req []domain.BatchPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, err := cookie.ReadEncrypted(c.Request, "id")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := h.URLsService.CreateBatch(c.Request.Context(), req, userID, h.cfg.BaseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)

}

func (h *Handler) Ping(c *gin.Context) {
	err := ping.Ping(c.Request.Context(), h.cfg.DatabaseDSN)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	}

	c.String(http.StatusOK, "")
}
