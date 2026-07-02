package models

import "time"

type User struct {
	ID           uint    `gorm:"primaryKey"`
	CompanyID    uint    `gorm:"not null;index"`
	Company      Company `gorm:"foreignKey:CompanyID"`
	Username     string  `gorm:"size:50;uniqueIndex;not null"`
	PasswordHash string  `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
