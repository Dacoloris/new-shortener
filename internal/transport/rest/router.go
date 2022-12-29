package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

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

func SetPlainTextHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "text/plain")
		c.Next()
	}
}

func SetJSONHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Next()
	}
}
