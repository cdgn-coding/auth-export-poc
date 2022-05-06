package repository

import (
	"auth0-users-sync-job-poc/internal/business/domain"
	"gorm.io/gorm"
)

type operatorDao struct {
	db *gorm.DB
}

func NewOperatorDao(db *gorm.DB) *operatorDao {
	return &operatorDao{db}
}

func (o operatorDao) Save(operator *domain.Operator) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		var err error
		err = o.db.Model(operator).Omit("Organizations").Save(operator).Error

		if err != nil {
			return err
		}

		if len(operator.Organizations) > 0 {
			err = o.db.Model(operator).Association("Organizations").Replace(operator.Organizations)
			return err
		}

		return nil
	})
}

func (o operatorDao) Get(ID string) (*domain.Operator, error) {
	var operator *domain.Operator
	result := o.db.Where("id = ?", ID).First(operator)
	return operator, result.Error
}
