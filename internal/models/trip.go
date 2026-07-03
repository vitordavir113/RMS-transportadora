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

	TractorID uint    `gorm:"not null;index" json:"tractor_id"`
	Tractor   Tractor `gorm:"foreignKey:TractorID" json:"tractor"`

	TrailerID uint    `gorm:"not null;index" json:"trailer_id"`
	Trailer   Trailer `gorm:"foreignKey:TrailerID" json:"trailer"`

	DriverID uint   `gorm:"not null;index" json:"driver_id"`
	Driver   Driver `gorm:"foreignKey:DriverID" json:"driver"`

	TractorPlateSnapshot string `gorm:"size:10;not null" json:"tractor_plate_snapshot"`
	TrailerPlateSnapshot string `gorm:"size:10;not null" json:"trailer_plate_snapshot"`
	DriverNameSnapshot   string `gorm:"size:150;not null" json:"driver_name_snapshot"`

	Status TripStatus `gorm:"size:20;not null;index;default:EM_ANDAMENTO" json:"status"`

	Compartments []TripCompartment `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE" json:"compartments"`

	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

type TripCompartment struct {
	ID uint `gorm:"primaryKey" json:"id"`

	TripID uint `gorm:"not null;index" json:"trip_id"`

	TrailerCompartmentID uint `gorm:"not null" json:"trailer_compartment_id"`

	ClientID uint   `gorm:"not null;index" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	Numero           int     `gorm:"not null" json:"numero"`
	CapacidadeLitros float64 `gorm:"not null" json:"capacidade_litros"`
	Produto          string  `gorm:"size:50;not null" json:"produto"`

	ClientNameSnapshot   string      `gorm:"size:150;not null" json:"client_name_snapshot"`
	FreightValueSnapshot float64     `gorm:"not null;default:0" json:"freight_value_snapshot"`
	FreightTypeSnapshot  FreightType `gorm:"size:20;not null" json:"freight_type_snapshot"`
	FreightTotal         float64     `gorm:"not null;default:0" json:"freight_total"`
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
