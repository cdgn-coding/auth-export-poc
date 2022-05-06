package config

import (
	"auth0-users-sync-job-poc/internal/infraestructure/env"
	"fmt"
	"github.com/spf13/viper"
	"path"
)

var (
	cfg *viper.Viper
)

func GetString(key string) string {
	return cfg.GetString(key)
}

func LoadConfig() {
	fname, filePath := getFileResource()
	cfg = viper.GetViper()
	cfg.SetConfigName(fname)
	cfg.AddConfigPath(filePath)
	cfg.AutomaticEnv()
	cfg.SetConfigType("yml")
	err := cfg.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func getFileResource() (string, string) {
	fname := env.GetEnvironment()

	if fname == "" {
		fname = "local"
	}

	basePath := env.GetWd()
	filePath := path.Join(basePath, "config")

	return fname, filePath
}
