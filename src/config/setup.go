package config

import (
	"fmt"
	"os"

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

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
