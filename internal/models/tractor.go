package models

import "time"

type Tractor struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID"`

	Model  string `gorm:"size:100;not null" json:"model"`
	Plate  string `gorm:"size:10;not null;index" json:"plate"`
	Owner  string `gorm:"size:150" json:"owner"`
	Status string `gorm:"size:30;not null;default:ativo" json:"status"`
	Active bool   `gorm:"not null;default:true" json:"active"`
	Notes  string `gorm:"type:text" json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
