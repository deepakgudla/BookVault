package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deepakgudla/BookVault/internal/config"
	"github.com/deepakgudla/BookVault/internal/database"
	"github.com/deepakgudla/BookVault/internal/interfaces"
	"github.com/deepakgudla/BookVault/internal/logger"
	"github.com/deepakgudla/BookVault/internal/providers"
	"github.com/deepakgudla/BookVault/internal/server"
	"github.com/deepakgudla/BookVault/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("database connection lost")
	}

	defer func() {
		if err := mainDB.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}()

	gin.SetMode(cfg.Server.GinMode)

	authService := services.NewAuthService(db, cfg)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)

	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	uploadService := services.NewUploadService(uploadProvider)

	serve := server.New(cfg, db, &log, authService, productService, userService, uploadService)

	router := serve.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  11 * time.Second,
		WriteTimeout: 11 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start the http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
		return
	}

	log.Info().Msg("shutting down database")
}
