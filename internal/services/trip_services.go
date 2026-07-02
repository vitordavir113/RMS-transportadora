package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"transportadora/internal/models"
)

type TripService struct {
	DB *gorm.DB
}

func NewTripService(db *gorm.DB) *TripService {
	return &TripService{DB: db}
}

type BocaInput struct {
	TruckCompartmentID uint
	Cliente            string
	Produto            string
}

func (s *TripService) CreateTrip(companyID uint, truckID uint, bocas []BocaInput) (*models.Trip, error) {
	if companyID == 0 {
		return nil, errors.New("empresa inválida")
	}

	if truckID == 0 {
		return nil, errors.New("selecione um caminhão")
	}

	var truck models.Truck
	if err := s.DB.Preload("Compartments").
		Where("company_id = ?", companyID).
		First(&truck, truckID).Error; err != nil {
		return nil, errors.New("caminhão não encontrado")
	}

	compartmentByID := map[uint]models.TruckCompartment{}
	for _, c := range truck.Compartments {
		compartmentByID[c.ID] = c
	}

	trip := models.Trip{
		CompanyID: companyID,
		TruckID:   truckID,
		Status:    models.StatusEmAndamento,
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trip).Error; err != nil {
			return err
		}

		for _, b := range bocas {
			if b.Cliente == "" || b.Produto == "" {
				continue
			}

			compartment, ok := compartmentByID[b.TruckCompartmentID]
			if !ok {
				return errors.New("boca inválida para este caminhão")
			}

			tripCompartment := models.TripCompartment{
				TripID:             trip.ID,
				TruckCompartmentID: compartment.ID,
				Numero:             compartment.Numero,
				CapacidadeLitros:   compartment.CapacidadeLitros,
				Cliente:            b.Cliente,
				Produto:            b.Produto,
			}

			if err := tx.Create(&tripCompartment).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.DB.Preload("Truck").Preload("Compartments").First(&trip, trip.ID)

	return &trip, nil
}

func (s *TripService) FinishTrip(companyID uint, tripID uint) error {
	now := time.Now()

	result := s.DB.Model(&models.Trip{}).
		Where("id = ? AND company_id = ? AND status = ?", tripID, companyID, models.StatusEmAndamento).
		Updates(map[string]interface{}{
			"status":      models.StatusFinalizada,
			"finished_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("viagem não encontrada ou já finalizada")
	}

	return nil
}
