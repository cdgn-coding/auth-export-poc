package usecase

import (
	"auth0-users-sync-job-poc/internal/business/domain"
	"auth0-users-sync-job-poc/internal/business/gateway"
	"auth0-users-sync-job-poc/internal/infraestructure/files"
	json "encoding/json"
	"github.com/auth0/go-auth0/management"
	"sync"
)

type UsersSyncer struct {
	auth0Service    gateway.Auth0Service
	operatorService gateway.OperatorService
}

func NewUsersSyncer(auth0Service gateway.Auth0Service, operatorService gateway.OperatorService) UsersSyncer {
	return UsersSyncer{
		auth0Service:    auth0Service,
		operatorService: operatorService,
	}
}

func (syncer UsersSyncer) Run() error {
	createdJob, errorRequestingExport := syncer.auth0Service.RequestExportUsers()
	if errorRequestingExport != nil {
		return errorRequestingExport
	}

	completedJob, errorDuringJobCompetition := syncer.auth0Service.WaitJobCompetition(*createdJob.ID)
	if errorDuringJobCompetition != nil {
		return errorDuringJobCompetition
	}

	filepath, errorGettingFile := syncer.auth0Service.GetUsersFile(*completedJob)
	if errorGettingFile != nil {
		return errorGettingFile
	}

	// Iterate File
	lineReader := files.NewReaderByLine()
	usersJson := make(chan string, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go syncer.updateOperators(usersJson, &wg)
	lineReader.Read(filepath, usersJson)

	wg.Wait()
	return nil
}

func (syncer *UsersSyncer) updateOperators(usersJson <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for userJson := range usersJson {
		var user management.User
		errorUnmarshalling := json.Unmarshal([]byte(userJson), &user)
		if errorUnmarshalling != nil {
			panic(errorUnmarshalling)
		}

		operator := domain.Operator{
			ID:            user.GetID(),
			Name:          user.GetName(),
			Email:         user.GetEmail(),
			Organizations: make([]domain.Organization, 0),
		}

		auth0OrgList, errorFetchingOrgs := syncer.auth0Service.GetUserOrganizations(user.GetID())
		if errorFetchingOrgs != nil {
			panic(errorFetchingOrgs)
		}
		for _, auth0Org := range auth0OrgList {
			operator.Organizations = append(operator.Organizations, domain.Organization{
				ID:          auth0Org.GetID(),
				Name:        auth0Org.GetName(),
				DisplayName: auth0Org.GetDisplayName(),
			})
		}
		errorSavingOperator := syncer.operatorService.Save(&operator)

		if errorSavingOperator != nil {
			panic(errorSavingOperator)
		}
	}
}
