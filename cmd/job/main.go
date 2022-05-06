package main

import (
	"auth0-users-sync-job-poc/internal/infraestructure/app"
)

func main() {
	job := app.NewAuth0SyncJob()
	job.Start()
}
