package config

type DbConfig struct {
	Host     string
	Database string
	User     string
	Port     uint16
	Password string
}

type Config struct {
	TargetDb DbConfig
}
