package app

import (
	"auth0-users-sync-job-poc/internal/business/external/auth0"
	"auth0-users-sync-job-poc/internal/business/repository"
	"auth0-users-sync-job-poc/internal/business/usecase"
	"auth0-users-sync-job-poc/internal/infraestructure/config"
	"auth0-users-sync-job-poc/internal/infraestructure/db"
)

type Auth0SyncJob struct {
	usersSyncer usecase.UsersSyncer
}

func NewAuth0SyncJob() *Auth0SyncJob {
	a := &Auth0SyncJob{}
	a.loadDependencies()
	return a
}

func (a *Auth0SyncJob) loadDependencies() {
	config.LoadConfig()

	auth0Client := auth0.NewClient()
	auth0Service := auth0.NewService(auth0Client)

	database := db.GetDBConnection()
	operatorService := repository.NewOperatorDao(database)
	usersSyncer := usecase.NewUsersSyncer(auth0Service, operatorService)
	a.usersSyncer = usersSyncer
}

func (a *Auth0SyncJob) Start() {
	a.usersSyncer.Run()
}
