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
	TrailerCompartmentID uint
	ClientID             uint
	Produto              string
}

func (s *TripService) CreateTrip(companyID uint, tractorID uint, trailerID uint, driverID uint, bocas []BocaInput) (*models.Trip, error) {
	if companyID == 0 {
		return nil, errors.New("empresa inválida")
	}

	if tractorID == 0 {
		return nil, errors.New("selecione um cavalo")
	}

	if trailerID == 0 {
		return nil, errors.New("selecione um tanque/carreta")
	}

	if driverID == 0 {
		return nil, errors.New("selecione um motorista")
	}

	var tractor models.Tractor
	if err := s.DB.Where("company_id = ? AND active = true", companyID).First(&tractor, tractorID).Error; err != nil {
		return nil, errors.New("cavalo não encontrado")
	}

	var trailer models.Trailer
	if err := s.DB.Preload("Compartments").
		Where("company_id = ? AND active = true", companyID).
		First(&trailer, trailerID).Error; err != nil {
		return nil, errors.New("tanque/carreta não encontrado")
	}

	var driver models.Driver
	if err := s.DB.Where("company_id = ? AND active = true", companyID).First(&driver, driverID).Error; err != nil {
		return nil, errors.New("motorista não encontrado")
	}

	compartmentByID := map[uint]models.TrailerCompartment{}
	for _, c := range trailer.Compartments {
		compartmentByID[c.ID] = c
	}

	trip := models.Trip{
		CompanyID:            companyID,
		TractorID:            tractor.ID,
		TrailerID:            trailer.ID,
		DriverID:             driver.ID,
		TractorPlateSnapshot: tractor.Plate,
		TrailerPlateSnapshot: trailer.Plate,
		DriverNameSnapshot:   driver.Name,
		Status:               models.StatusEmAndamento,
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trip).Error; err != nil {
			return err
		}

		for _, b := range bocas {
			if b.ClientID == 0 || b.Produto == "" {
				continue
			}

			compartment, ok := compartmentByID[b.TrailerCompartmentID]
			if !ok {
				return errors.New("boca inválida para este tanque/carreta")
			}

			var client models.Client
			if err := tx.Where("company_id = ? AND active = true", companyID).First(&client, b.ClientID).Error; err != nil {
				return errors.New("cliente não encontrado")
			}

			freightTotal := calculateFreight(client.FreightType, client.FreightValue, compartment.CapacidadeLitros)

			tripCompartment := models.TripCompartment{
				TripID:               trip.ID,
				TrailerCompartmentID: compartment.ID,
				ClientID:             client.ID,
				Numero:               compartment.Numero,
				CapacidadeLitros:     compartment.CapacidadeLitros,
				Produto:              b.Produto,
				ClientNameSnapshot:   client.Name,
				FreightValueSnapshot: client.FreightValue,
				FreightTypeSnapshot:  client.FreightType,
				FreightTotal:         freightTotal,
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

	s.DB.
		Preload("Tractor").
		Preload("Trailer").
		Preload("Driver").
		Preload("Compartments").
		First(&trip, trip.ID)

	return &trip, nil
}

func calculateFreight(freightType models.FreightType, value float64, liters float64) float64 {
	switch freightType {
	case models.FreightPorLitro:
		return value * liters
	case models.FreightPorViagem:
		return value
	case models.FreightPorBoca:
		return value
	default:
		return value
	}
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
