package rest

import (
	"new-shortner/internal/transport/rest/cookie"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cookie.CheckCookie(), adapter.Wrap(GzipInput), adapter.Wrap(GzipOutput))

	plainText := r.Group("/")
	{
		plainText.Use(SetPlainTextHeader())
		plainText.GET("/:id", h.Redirect)
		plainText.POST("/", h.URLShortening)
		plainText.GET("/ping", h.Ping)
	}

	api := r.Group("/api")
	{
		api.Use(SetJSONHeader())
		api.POST("/shorten", h.APIShorten)
		api.GET("/user/urls", h.GetAllURLsForUser)
	}

	return r
}
