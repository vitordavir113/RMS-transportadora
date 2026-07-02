package models

import "time"

type Company struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:120;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
