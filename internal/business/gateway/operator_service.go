package gateway

import (
	"auth0-users-sync-job-poc/internal/business/domain"
	"time"
)

type OperatorService interface {
	Save(operator *domain.Operator) error
	Get(ID string) (*domain.Operator, error)
	OuterDelete(IDs []string, timeLimit time.Time) error
}
