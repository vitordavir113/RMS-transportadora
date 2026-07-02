package models

import "time"

type TripStatus string

const (
	StatusEmAndamento TripStatus = "EM_ANDAMENTO"
	StatusFinalizada  TripStatus = "FINALIZADA"
)

// Produtos fixos do sistema — não possuem cadastro em banco.
var ProdutosDisponiveis = []string{
	"Gasolina Comum",
	"Gasolina Aditivada",
	"Diesel S10",
	"Diesel S500",
	"Etanol",
}

// Trip representa uma viagem de um caminhão, com uma ou mais bocas carregadas.
type Trip struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	CompanyID    uint              `gorm:"not null;index" json:"company_id"`
	Company      Company           `gorm:"foreignKey:CompanyID" json:"company"`
	TruckID      uint              `gorm:"not null;index" json:"truck_id"`
	Truck        Truck             `gorm:"foreignKey:TruckID" json:"truck"`
	Status       TripStatus        `gorm:"size:20;not null;index;default:EM_ANDAMENTO" json:"status"`
	Compartments []TripCompartment `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE" json:"compartments"`
	CreatedAt    time.Time         `json:"created_at"`
	FinishedAt   *time.Time        `json:"finished_at"`
}

// TripCompartment é a boca de uma viagem específica: cliente + produto,
// com número e capacidade copiados da boca do caminhão no momento da criação.
type TripCompartment struct {
	ID                 uint    `gorm:"primaryKey" json:"id"`
	TripID             uint    `gorm:"not null;index" json:"trip_id"`
	TruckCompartmentID uint    `gorm:"not null" json:"truck_compartment_id"`
	Numero             int     `gorm:"not null" json:"numero"`
	CapacidadeLitros   float64 `gorm:"not null" json:"capacidade_litros"`
	Cliente            string  `gorm:"size:150;not null" json:"cliente"`
	Produto            string  `gorm:"size:50;not null" json:"produto"`
}

// TotalLitros soma a capacidade de todas as bocas da viagem.
func (t *Trip) TotalLitros() float64 {
	var total float64
	for _, c := range t.Compartments {
		total += c.CapacidadeLitros
	}
	return total
}

// QuantidadeClientes conta clientes distintos na viagem.
func (t *Trip) QuantidadeClientes() int {
	seen := map[string]bool{}
	for _, c := range t.Compartments {
		if c.Cliente != "" {
			seen[c.Cliente] = true
		}
	}
	return len(seen)
}
