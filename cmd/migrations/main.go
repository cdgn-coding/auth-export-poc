package main

import (
	"auth0-users-sync-job-poc/internal/business/domain"
	"auth0-users-sync-job-poc/internal/infraestructure/config"
	"auth0-users-sync-job-poc/internal/infraestructure/db"
	"fmt"
)

func main() {
	config.LoadConfig()
	database := db.GetDBConnection()

	err := database.AutoMigrate(&domain.Organization{}, &domain.Operator{})
	if err != nil {
		panic(fmt.Errorf("error during migration %v", err))
	}
}
