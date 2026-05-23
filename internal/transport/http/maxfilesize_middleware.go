package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LimitBodySize(maxMb int64) gin.HandlerFunc {
	maxBytes := maxMb * 1024 * 1024

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		c.Next()
		if len(c.Errors) > 0 {
			var maxBytesErr *http.MaxBytesError
			if errors.As(c.Errors.Last().Err, &maxBytesErr) {
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
					"error":       "file too large",
					"max_size_mb": maxMb,
				})
			}
		}
	}
}
