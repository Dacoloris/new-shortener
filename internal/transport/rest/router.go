package rest

import (
	"net/http"
	"new-shortner/internal/transport/rest/cookie"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(adapter.Wrap(GzipInput), adapter.Wrap(GzipOutput), cookie.CheckCookie)

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
		api.GET("/user/urls", h.GetAllURLsForUser)
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	return r
}
