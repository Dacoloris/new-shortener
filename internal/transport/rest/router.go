package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	//r.Use(GzipInput(), adapter.Wrap(GzipOutput))

	plainText := r.Group("/")
	{
		plainText.Use(SetPlainTextHeader())
		plainText.GET("/:id", h.Redirect)
		plainText.POST("/", h.URLShortening)
	}

	api := r.Group("/api")
	{
		api.Use(SetJSONHeader())
		api.POST("/shorten", h.APIShorten)
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	return r
}
