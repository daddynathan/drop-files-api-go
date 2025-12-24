package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LimitBodySize(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitedBody := http.MaxBytesReader(c.Writer, c.Request.Body, max)
		var buf bytes.Buffer
		tmp := make([]byte, 32*1024) // 32KB буфер
		for {
			n, err := limitedBody.Read(tmp)
			if n > 0 {
				if _, werr := buf.Write(tmp[:n]); werr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "failed to read request body",
					})
					return
				}
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				var tooLarge *http.MaxBytesError
				if errors.As(err, &tooLarge) {
					c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
						"error":       "file too large",
						"max_size_mb": (max / (1024 * 1024)) - 1,
					})
					return
				}
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "invalid request body",
				})
				return
			}
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		c.Next()
	}
}
