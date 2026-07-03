package models

import "time"

type TripStatus string

const (
	StatusEmAndamento TripStatus = "EM_ANDAMENTO"
	StatusFinalizada  TripStatus = "FINALIZADA"
)

var ProdutosDisponiveis = []string{
	"Gasolina Comum",
	"Gasolina Aditivada",
	"Diesel S10",
	"Diesel S500",
	"Etanol",
}

type Trip struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID" json:"company"`

	TractorID uint `gorm:"index"`
	TrailerID uint `gorm:"index"`
	DriverID  uint `gorm:"index"`

	TractorPlateSnapshot string `gorm:"size:10"`
	TrailerPlateSnapshot string `gorm:"size:10"`
	DriverNameSnapshot   string `gorm:"size:150"`

	Status TripStatus `gorm:"size:20;not null;index;default:EM_ANDAMENTO" json:"status"`

	Compartments []TripCompartment `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE" json:"compartments"`

	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

type TripCompartment struct {
	ID uint `gorm:"primaryKey" json:"id"`

	TripID               uint `gorm:"not null;index" json:"trip_id"`
	TrailerCompartmentID uint
	ClientID             uint `gorm:"index"`

	ClientNameSnapshot   string      `gorm:"size:150"`
	FreightValueSnapshot float64     `gorm:"default:0"`
	FreightTypeSnapshot  FreightType `gorm:"size:20"`
	FreightTotal         float64     `gorm:"default:0"`
	Client               Client      `gorm:"foreignKey:ClientID" json:"client"`

	Numero           int     `gorm:"not null" json:"numero"`
	CapacidadeLitros float64 `gorm:"not null" json:"capacidade_litros"`
	Produto          string  `gorm:"size:50;not null" json:"produto"`
}

func (t Trip) TotalFrete() float64 {
	total := 0.0

	for _, c := range t.Compartments {
		total += c.FreightTotal
	}

	return total
}

func (t Trip) TotalLitros() float64 {
	total := 0.0

	for _, c := range t.Compartments {
		total += c.CapacidadeLitros
	}

	return total
}

func (t *Trip) QuantidadeClientes() int {
	seen := map[string]bool{}
	for _, c := range t.Compartments {
		if c.ClientNameSnapshot != "" {
			seen[c.ClientNameSnapshot] = true
		}
	}
	return len(seen)
}
