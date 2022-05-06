package usecase

import (
	"auth0-users-sync-job-poc/internal/business/domain"
	"auth0-users-sync-job-poc/internal/business/gateway"
	"auth0-users-sync-job-poc/internal/infraestructure/files"
	json "encoding/json"
	"github.com/auth0/go-auth0/management"
	"sync"
	"time"
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

	jobCreationTime := *completedJob.CreatedAt
	syncer.syncWithFile(jobCreationTime, filepath)

	return nil
}

func (syncer *UsersSyncer) syncWithFile(jobCreationTime time.Time, filepath string) {
	lineReader := files.NewReaderByLine()
	jsonLine := make(chan string, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go syncer.processLineByLine(jobCreationTime, jsonLine, &wg)

	wg.Add(1)
	go lineReader.Read(filepath, jsonLine, &wg)

	wg.Wait()
}

func (syncer *UsersSyncer) processLineByLine(jobCreationTime time.Time, usersJson <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var existingUsers = make([]string, 0)

	for userJson := range usersJson {
		operator := syncer.updateOperator(userJson)
		existingUsers = append(existingUsers, operator.ID)
	}

	err := syncer.operatorService.OuterDelete(existingUsers, jobCreationTime)
	if err != nil {
		return
	}
}

func (syncer *UsersSyncer) updateOperator(userJson string) *domain.Operator {
	var user management.User
	var err error
	var operator *domain.Operator
	err = json.Unmarshal([]byte(userJson), &user)

	if err != nil {
		panic(err)
	}

	operator, err = syncer.toDomainOperator(user)

	err = syncer.operatorService.Save(operator)

	if err != nil {
		panic(err)
	}
	return operator
}

func (syncer *UsersSyncer) toDomainOperator(user management.User) (*domain.Operator, error) {
	operator := &domain.Operator{
		ID:    user.GetID(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}

	orgList, err := syncer.auth0Service.GetUserOrganizations(user.GetID())
	if err != nil {
		return nil, err
	}

	operator.Organizations = make([]domain.Organization, len(orgList))
	for i, auth0Org := range orgList {
		operator.Organizations[i] = syncer.toDomainOrganization(auth0Org)
	}

	return operator, nil
}

func (syncer *UsersSyncer) toDomainOrganization(auth0Org *management.Organization) domain.Organization {
	return domain.Organization{
		ID:          auth0Org.GetID(),
		Name:        auth0Org.GetName(),
		DisplayName: auth0Org.GetDisplayName(),
	}
}
