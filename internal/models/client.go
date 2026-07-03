package models

import "time"

type FreightType string

const (
	FreightPorBoca   FreightType = "POR_BOCA"
	FreightPorLitro  FreightType = "POR_LITRO"
	FreightPorViagem FreightType = "POR_VIAGEM"
)

type Client struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID"`

	Name         string      `gorm:"size:150;not null;index" json:"name"`
	Document     string      `gorm:"size:30" json:"document"`
	Phone        string      `gorm:"size:30" json:"phone"`
	FreightValue float64     `gorm:"not null;default:0" json:"freight_value"`
	FreightType  FreightType `gorm:"size:20;not null;default:POR_BOCA" json:"freight_type"`
	Active       bool        `gorm:"not null;default:true" json:"active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
