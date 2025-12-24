package http

import (
	"DropFiles/config"
	"DropFiles/internal/file"
	"context"
	"log"
	"net/http"
	"time"

	// _ "DropFiles/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HTTPServer struct {
	server *http.Server
}

type HTTPHandlers struct {
	UploadHandler   gin.HandlerFunc
	DownloadHandler gin.HandlerFunc
}

func NewHTTPHandlers(fileHandler *file.Handler) *HTTPHandlers {
	return &HTTPHandlers{
		UploadHandler:   fileHandler.UploadHandler,
		DownloadHandler: fileHandler.DownloadHandler,
	}
}

func NewHTTPServer(handlers *HTTPHandlers, addr string) *HTTPServer {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	apiGroup := router.Group("/api")
	{
		v1 := apiGroup.Group("/v1")
		v1.Use(LimitBodySize(config.Load().File.MaxFileSizeBytes + 1<<20)) // +1MB
		{
			v1.GET("/f/:id", handlers.DownloadHandler)
			v1.POST("/upload", handlers.UploadHandler)
		}
	}
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return &HTTPServer{
		server: httpServer,
	}
}

func (s *HTTPServer) Start() {
	log.Printf("HTTP server starting on %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Addr() string {
	return s.server.Addr
}
