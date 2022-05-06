package domain

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	ID          string `gorm:"primary_key,not null"`
	Name        string `gorm:"not null"`
	DisplayName string `gorm:"not null"`
}
