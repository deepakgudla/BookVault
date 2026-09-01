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

	"github.com/deepakgudla/bookvault/internal/config"
	"github.com/deepakgudla/bookvault/internal/database"
	"github.com/deepakgudla/bookvault/internal/events"
	"github.com/deepakgudla/bookvault/internal/interfaces"
	"github.com/deepakgudla/bookvault/internal/logger"
	"github.com/deepakgudla/bookvault/internal/providers"
	"github.com/deepakgudla/bookvault/internal/repository"
	"github.com/deepakgudla/bookvault/internal/server"
	"github.com/deepakgudla/bookvault/internal/services"
	"github.com/gin-gonic/gin"

	_ "github.com/deepakgudla/bookvault/docs"
)

// @title bookvault-API
// @version 1.0
// @description A modern e-commerce API built with Go
// @termsOfService http://swagger.io/terms/

// @contact.name Deepak Gudla
// @contact.url https://www.linkedin.com/in/deepakgudla
// @contact.email deepakgudla35@protonmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1357
// @BasePath /api
// @schemas http https

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Bearer and JWT token (bearer jwt_token)

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

	ctx := context.Background()

	eventPublisher, err := events.NewEventPublisher(ctx, &cfg.AWS)
	if err != nil {
		log.Error().Err(err).Msg("failed to create event publisher")
		return
	}

	gin.SetMode(cfg.Server.GinMode)

	userRepository := repository.NewUserRepository(db)
	cartRepository := repository.NewCartRepository(db)

	authService := services.NewAuthService(cfg, eventPublisher, userRepository, cartRepository)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)
	cartService := services.NewCartService(db)
	orderService := services.NewOrderService(db)

	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	uploadService := services.NewUploadService(uploadProvider)

	serve := server.New(
		cfg,
		// db,
		&log,
		authService,
		productService,
		userService,
		uploadService,
		cartService,
		orderService,
	)

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
