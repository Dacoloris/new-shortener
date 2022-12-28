package rest

import (
	"context"
	"net/http"
	"new-shortner/internal/config"
	"new-shortner/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Urls interface {
	Create(ctx context.Context, book domain.URL) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.URL, error)
	GetAll(ctx context.Context) ([]domain.URL, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ShortenUrl(ctx context.Context) string
}

type Handler struct {
	UrlsService Urls
	cfg         config.Config
}

func NewHandler(urls Urls, cfg config.Config) *Handler {
	return &Handler{
		UrlsService: urls,
		cfg:         cfg,
	}
}

func (h *Handler) Redirect(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
	}

	url, err := h.UrlsService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
	}

	c.Header("Location", url.Original)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) UrlShortening(c *gin.Context) {
	var original string
	err := c.ShouldBind(&original)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	short := h.UrlsService.ShortenUrl(c.Request.Context())
	c.String(http.StatusCreated, "%s/%s", h.cfg.BaseURL, short)
}
