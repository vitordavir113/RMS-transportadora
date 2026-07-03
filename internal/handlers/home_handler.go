package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type HomeHandler struct {
	DB *gorm.DB
}

func NewHomeHandler(db *gorm.DB) *HomeHandler {
	return &HomeHandler{DB: db}
}

type DashboardTripView struct {
	models.Trip
	TotalFrete  float64
	TotalLitros float64
}

func (h *HomeHandler) Index(c *gin.Context) {
	companyID := CurrentCompanyID(c)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	var trips []models.Trip
	h.DB.
		Preload("Compartments").
		Where("company_id = ? AND status = ?", companyID, models.StatusEmAndamento).
		Order("created_at desc").
		Limit(6).
		Find(&trips)

	var viewTrips []DashboardTripView
	var activeLiters float64
	var activeFreight float64

	for _, trip := range trips {
		item := DashboardTripView{Trip: trip}

		for _, b := range trip.Compartments {
			item.TotalLitros += b.CapacidadeLitros
			item.TotalFrete += b.FreightTotal
		}

		activeLiters += item.TotalLitros
		activeFreight += item.TotalFrete
		viewTrips = append(viewTrips, item)
	}

	var tractorsCount int64
	h.DB.Model(&models.Tractor{}).
		Where("company_id = ? AND active = true", companyID).
		Count(&tractorsCount)

	var trailersCount int64
	h.DB.Model(&models.Trailer{}).
		Where("company_id = ? AND active = true", companyID).
		Count(&trailersCount)

	var clientsCount int64
	h.DB.Model(&models.Client{}).
		Where("company_id = ? AND active = true", companyID).
		Count(&clientsCount)

	var monthlyFreight float64
	h.DB.Model(&models.TripCompartment{}).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trips.company_id = ? AND trips.created_at >= ? AND trips.created_at < ?", companyID, monthStart, monthEnd).
		Select("COALESCE(SUM(trip_compartments.freight_total), 0)").
		Scan(&monthlyFreight)

	var finishedThisMonth int64
	h.DB.Model(&models.Trip{}).
		Where("company_id = ? AND status = ? AND created_at >= ? AND created_at < ?", companyID, models.StatusFinalizada, monthStart, monthEnd).
		Count(&finishedThisMonth)

	RenderPage(c, "home/index.html", gin.H{
		"Title":             "Dashboard",
		"Trips":             viewTrips,
		"ActiveTrips":       len(viewTrips),
		"ActiveLiters":      activeLiters,
		"ActiveFreight":     activeFreight,
		"MonthlyFreight":    monthlyFreight,
		"FinishedThisMonth": finishedThisMonth,
		"TractorsCount":     tractorsCount,
		"TrailersCount":     trailersCount,
		"ClientsCount":      clientsCount,
	})
}
