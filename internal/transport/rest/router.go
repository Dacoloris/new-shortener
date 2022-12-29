package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(SetPlainTextHeader())
	r.GET("/:id", h.Redirect)
	r.POST("/", h.URLShortening)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func SetPlainTextHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "text/plain")
		c.Next()
	}
}
