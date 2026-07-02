package models

import "time"

// Truck representa o cavalo mecânico + tanque cadastrado, compania e usuario

type Truck struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID"`

	CavaloModelo string `gorm:"size:100;not null" json:"cavalo_modelo"`
	CavaloPlaca  string `gorm:"size:10;not null;index" json:"cavalo_placa"`
	TanquePlaca  string `gorm:"size:10;not null" json:"tanque_placa"`

	Compartments []TruckCompartment `gorm:"foreignKey:TruckID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TruckCompartment representa uma "boca" do tanque, com capacidade fixa.
type TruckCompartment struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TruckID          uint      `gorm:"not null;index" json:"truck_id"`
	Numero           int       `gorm:"not null" json:"numero"`
	CapacidadeLitros float64   `gorm:"not null" json:"capacidade_litros"`
	CreatedAt        time.Time `json:"created_at"`
}

// TotalCapacidade soma a capacidade de todas as bocas do caminhão.
func (t *Truck) TotalCapacidade() float64 {
	var total float64
	for _, c := range t.Compartments {
		total += c.CapacidadeLitros
	}
	return total
}
