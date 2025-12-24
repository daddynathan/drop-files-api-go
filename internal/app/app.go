package app

import (
	"errors"
	"log"
	"os"

	"DropFiles/config"
	"DropFiles/internal/file"
	"DropFiles/internal/infra/db"
	httpD "DropFiles/internal/transport/http"

	"github.com/joho/godotenv"
)

type App struct {
	HTTPServer *httpD.HTTPServer
}

func New() (*App, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("using system environment variables")
	}
	cfg := config.Load()
	if cfg.Server.Port == "" {
		return nil, errors.New("APP_PORT (SERVER_PORT) is required")
	}
	if err := os.MkdirAll(cfg.File.UploadDir, 0755); err != nil {
		return nil, err
	}
	sqlDB, err := db.NewPostgresDB()
	if err != nil {
		return nil, err
	}
	fileRepo := file.NewFileRepo(sqlDB)
	fileStorage := file.NewLocalStorage(cfg.File.UploadDir)
	fileService := file.NewFileService(fileRepo, fileStorage, cfg.File)
	fileHandler := file.NewHandler(fileService, cfg.File.UploadDir)

	httpHandlers := httpD.NewHTTPHandlers(fileHandler)
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	httpServer := httpD.NewHTTPServer(httpHandlers, addr)

	return &App{
		HTTPServer: httpServer,
	}, nil
}

func (a *App) Run() {
	log.Printf("starting server on %s", a.HTTPServer.Addr())
	a.HTTPServer.Start()
}
