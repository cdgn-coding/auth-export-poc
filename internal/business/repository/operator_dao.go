package repository

import (
	"auth0-users-sync-job-poc/internal/business/domain"
	"gorm.io/gorm"
	"time"
)

type operatorDao struct {
	db *gorm.DB
}

func (o operatorDao) OuterDelete(IDs []string, timeLimit time.Time) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec("DELETE FROM operators WHERE (id NOT IN ?) AND (created_at < ?)", IDs, timeLimit)
		return result.Error
	})
}

func NewOperatorDao(db *gorm.DB) *operatorDao {
	return &operatorDao{db}
}

func (o operatorDao) Save(operator *domain.Operator) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if o.db.Model(&domain.Operator{}).Where("id = ?", operator.ID).Updates(&operator).RowsAffected == 0 {
			o.db.Create(&operator)
		}

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
