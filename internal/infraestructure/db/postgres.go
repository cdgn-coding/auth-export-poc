package db

import (
	"auth0-users-sync-job-poc/internal/infraestructure/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const postgresConfig = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

type Config struct {
	Host     string
	User     string
	Port     int
	Password string
	DBName   string
}

func getDBConfiguration() Config {
	return Config{
		Host:     config.GetString("database.host"),
		User:     config.GetString("database.user"),
		Port:     5432,
		Password: config.GetString("database.password"),
		DBName:   config.GetString("database.db"),
	}
}

func GetDBConnection() *gorm.DB {
	cfg := getDBConfiguration()
	dsn := fmt.Sprintf(postgresConfig, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
