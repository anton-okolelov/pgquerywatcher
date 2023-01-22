package main

import (
	"os"
	"strconv"

	"github.com/anton.okolelov/pgquerywatcher/internal/config"
	"github.com/anton.okolelov/pgquerywatcher/internal/watch"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout)
	log.Info().Msg("Start watching")

	port, err := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 16)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse DB_PORT")
	}

	cfg := config.Config{
		TargetDb: config.DbConfig{
			Host:     os.Getenv("DB_HOST"),
			Database: os.Getenv("DB_DATABASE"),
			User:     os.Getenv("DB_USER"),
			Port:     uint16(port),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}

	watcher, err := watch.NewWatcher(cfg.TargetDb, log)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create watcher")
	}

	err = watcher.Watch()
	if err != nil {
		log.Fatal().Err(err).Msg("error watching")
	}
}
