package auth0

import (
	"auth0-users-sync-job-poc/internal/infraestructure/config"
	"fmt"
	"github.com/auth0/go-auth0/management"
)

func NewClient() *management.Management {
	id := config.GetString("auth0.client_id")
	secret := config.GetString("auth0.client_secret")
	domain := config.GetString("auth0.domain")

	client, errorRequestingToken := management.New(domain, management.WithClientCredentials(id, secret))

	if errorRequestingToken != nil {
		panic(fmt.Errorf("cannot connect to Auth0 %v", errorRequestingToken))
	}

	return client
}
