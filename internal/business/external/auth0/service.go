package auth0

import (
	"auth0-users-sync-job-poc/internal/infraestructure/config"
	"auth0-users-sync-job-poc/internal/infraestructure/files"
	"fmt"
	"github.com/auth0/go-auth0/management"
	"path"
	"time"
)

type Service struct {
	auth0Client *management.Management
}

func (e Service) GetUserOrganizations(ID string) ([]*management.Organization, error) {
	organizationList, errorGettingOrganizations := e.auth0Client.User.Organizations(ID)

	if errorGettingOrganizations != nil {
		return nil, errorGettingOrganizations
	}

	return organizationList.Organizations, nil
}

func NewService(client *management.Management) Service {
	return Service{
		auth0Client: client,
	}
}

func (e Service) RequestExportUsers() (*management.Job, error) {
	connectionId := config.GetString("auth0.databaseConnection")
	exportJob := e.createExportUsersJob(connectionId)
	errorCreatingJob := e.auth0Client.Job.ExportUsers(exportJob)

	if errorCreatingJob != nil {
		return nil, errorCreatingJob
	}

	return exportJob, nil
}

func (e Service) WaitJobCompetition(jobID string) (*management.Job, error) {
	for {
		job, errorReadingJob := e.auth0Client.Job.Read(jobID)

		if errorReadingJob != nil {
			return nil, errorReadingJob
		}

		if *job.Status == "completed" {
			return job, nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (e Service) GetUsersFile(job management.Job) (string, error) {
	downloader := files.NewDownloader()
	filepath := path.Join("tmp", fmt.Sprintf("%s.csv", *job.ID))
	errorDownloading := downloader.Download(*job.Location, filepath)

	if errorDownloading != nil {
		return "", errorDownloading
	}

	return filepath, nil
}

func (e Service) createExportUsersJob(sourceDatabase string) *management.Job {
	format := "json"
	return &management.Job{
		ConnectionID: &sourceDatabase,
		Format:       &format,
		Fields: []map[string]interface{}{
			{
				"name": "user_id",
			},
			{
				"name": "name",
			},
			{
				"name": "email",
			},
			{
				"name": "user_metadata",
			},
			{
				"name": "app_metadata",
			},
		},
	}
}
