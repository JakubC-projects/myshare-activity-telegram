package config

import (
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/kelseyhightower/envconfig"
)

var cfg Config
var isLoaded = false

func Get() Config {
	if !isLoaded {
		readEnv(&cfg)
	}
	return cfg
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.L.Fatal().AnErr("err", err).Msg("Cannot read env")
	}
}
