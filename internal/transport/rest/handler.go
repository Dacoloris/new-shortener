package rest

import (
	"context"
	"net/http"
	"new-shortner/internal/config"

	"github.com/gin-gonic/gin"
)

type URLs interface {
	Create(ctx context.Context, original string) (string, error)
	GetOriginalByShort(ctx context.Context, short string) (string, error)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	c.Header("Location", original)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) URLShortening(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	original := string(b)

	short, _ := h.URLsService.Create(c.Request.Context(), original)

	c.String(http.StatusCreated, "%s/%s", h.cfg.BaseURL, short)
}
