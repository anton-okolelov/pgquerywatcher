package pg

import (
	"context"

	"github.com/anton.okolelov/pgquerywatcher/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDB(dbCfg config.DbConfig) (*pgxpool.Pool, error) {
	cfg, _ := pgxpool.ParseConfig("")
	cfg.ConnConfig.Host = dbCfg.Host
	cfg.ConnConfig.Port = dbCfg.Port
	cfg.ConnConfig.User = dbCfg.User
	cfg.ConnConfig.Password = dbCfg.Password
	cfg.ConnConfig.Database = dbCfg.Database
	cfg.ConnConfig.PreferSimpleProtocol = true

	dbPool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	return dbPool, err
}
