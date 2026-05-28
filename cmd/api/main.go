package main

import (
	"github.com/deepakgudla/BookVault/internal/config"
	"github.com/deepakgudla/BookVault/internal/database"
	"github.com/deepakgudla/BookVault/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("database connection lost")
	}

	defer mainDB.Close()
	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("server running")
}
