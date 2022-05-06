package domain

import (
	"gorm.io/gorm"
	"time"
)

type Operator struct {
	gorm.Model
	ID            string         `gorm:"primaryKey;not null"`
	Name          string         `gorm:"not null"`
	Email         string         `gorm:"index;not null"`
	Organizations []Organization `gorm:"many2many:operator_organizations;"`
	UpdatedAt     time.Time
}
