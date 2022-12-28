package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/:id", h.Redirect)
	r.POST("/", h.URLShortening)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
