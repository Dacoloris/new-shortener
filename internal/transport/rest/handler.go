package rest

import (
	"context"
	_ "encoding/json"
	"errors"
	"fmt"
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
		c.String(http.StatusBadRequest, errors.New("invalid id").Error())

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

	original := string(b)

	short, err := h.URLsService.Create(c.Request.Context(), original)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusCreated, "%s/%s", h.cfg.BaseURL, short)
}

func (h *Handler) APIShorten(c *gin.Context) {
	j := struct {
		URL string `json:"url"`
	}{}
	if err := c.ShouldBindJSON(&j); err != nil || j.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": errors.New("invalid json").Error()})
		return
	}

	short, err := h.URLsService.Create(c.Request.Context(), j.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": errors.New("invalid url").Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": fmt.Sprintf(`%s/%s`, h.cfg.BaseURL, short)})
}
