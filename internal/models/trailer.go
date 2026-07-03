package models

import "time"

type Trailer struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID"`

	Plate  string `gorm:"size:10;not null;index" json:"plate"`
	Type   string `gorm:"size:50;not null;default:tanque" json:"type"`
	Owner  string `gorm:"size:150" json:"owner"`
	Status string `gorm:"size:30;not null;default:ativo" json:"status"`
	Active bool   `gorm:"not null;default:true" json:"active"`
	Notes  string `gorm:"type:text" json:"notes"`

	Compartments []TrailerCompartment `gorm:"foreignKey:TrailerID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrailerCompartment struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	TrailerID        uint    `gorm:"not null;index" json:"trailer_id"`
	Numero           int     `gorm:"not null" json:"numero"`
	CapacidadeLitros float64 `gorm:"not null" json:"capacidade_litros"`

	CreatedAt time.Time `json:"created_at"`
}

func (t *Trailer) TotalCapacidade() float64 {
	var total float64
	for _, c := range t.Compartments {
		total += c.CapacidadeLitros
	}
	return total
}
