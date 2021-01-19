package main

import (
	"os"

	"github.com/anton.okolelov/pgquerywatcher/internal/config"
	"github.com/anton.okolelov/pgquerywatcher/internal/watch"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	log := zerolog.New(os.Stdout)
	log.Info().Msg("Start watching")
	targetDbPassword := os.Getenv("TARGET_DB_PASSWORD")

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "watch",
				Usage: "watch database queries. Password must be in $TARGET_DB_PASSWORD env var",
				Action: func(c *cli.Context) error {
					cfg := config.Config{
						TargetDb: config.DbConfig{
							Host:     c.String("target_db_host"),
							Database: c.String("target_db_database"),
							User:     c.String("target_db_user"),
							Port:     uint16(c.Uint("target_db_port")),
							Password: targetDbPassword,
						},
					}

					watcher, err := watch.NewWatcher(cfg.TargetDb, log)
					if err != nil {
						return err
					}
					err = watcher.Watch()
					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "target_db_host", Required: true, Usage: "watched pg host (required)"},
					&cli.StringFlag{Name: "target_db_database", Required: true, Usage: "watched pg database name (required)"},
					&cli.StringFlag{Name: "target_db_user", Required: true, Usage: "watched pg user (required)"},
					&cli.UintFlag{Name: "target_db_port", Required: true, Usage: "watched pg port (required)"},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error().Err(err)
	}
}
