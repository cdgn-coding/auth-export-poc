package gateway

import "github.com/auth0/go-auth0/management"

type Auth0Service interface {
	RequestExportUsers() (*management.Job, error)
	WaitJobCompetition(jobID string) (*management.Job, error)
	GetUsersFile(job management.Job) (string, error)
	GetUserOrganizations(ID string) ([]*management.Organization, error)
}
