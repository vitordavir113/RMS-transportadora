package models

import "time"

type Driver struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID"`

	Name   string `gorm:"size:150;not null" json:"name"`
	CPF    string `gorm:"size:20" json:"cpf"`
	Phone  string `gorm:"size:30" json:"phone"`
	CNH    string `gorm:"size:30" json:"cnh"`
	Active bool   `gorm:"not null;default:true" json:"active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
