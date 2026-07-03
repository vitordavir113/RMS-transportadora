package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type ReportHandler struct {
	DB *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{DB: db}
}

type ReportSummary struct {
	PeriodLabel      string
	TotalFreight     float64
	PreviousFreight  float64
	FreightVariation float64
	TotalLiters      float64
	FinishedTrips    int64
	ActiveTrips      int64
	AveragePerTrip   float64
	FreightByClient  []ClientFreightReport
}

type ClientFreightReport struct {
	ClientName string
	Total      float64
	Liters     float64
	Trips      int64
}

func (h *ReportHandler) Monthly(c *gin.Context) {
	companyID := CurrentCompanyID(c)

	now := time.Now()
	year := now.Year()
	month := now.Month()

	start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)

	previousStart := start.AddDate(0, -1, 0)
	previousEnd := start

	summary := h.buildReport(companyID, start, end, previousStart, previousEnd)

	RenderPage(c, "reports/monthly.html", gin.H{
		"Title":   "Relatório mensal",
		"Summary": summary,
	})
}

func (h *ReportHandler) Weekly(c *gin.Context) {
	companyID := CurrentCompanyID(c)

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -(weekday - 1))

	end := start.AddDate(0, 0, 7)

	previousStart := start.AddDate(0, 0, -7)
	previousEnd := start

	summary := h.buildReport(companyID, start, end, previousStart, previousEnd)

	RenderPage(c, "reports/monthly.html", gin.H{
		"Title":   "Relatório semanal",
		"Summary": summary,
	})
}

func (h *ReportHandler) buildReport(companyID uint, start, end, previousStart, previousEnd time.Time) ReportSummary {
	var totalFreight float64
	var previousFreight float64
	var totalLiters float64
	var finishedTrips int64
	var activeTrips int64

	h.DB.Model(&models.TripCompartment{}).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trips.company_id = ? AND trips.created_at >= ? AND trips.created_at < ?", companyID, start, end).
		Select("COALESCE(SUM(trip_compartments.freight_total), 0)").
		Scan(&totalFreight)

	h.DB.Model(&models.TripCompartment{}).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trips.company_id = ? AND trips.created_at >= ? AND trips.created_at < ?", companyID, previousStart, previousEnd).
		Select("COALESCE(SUM(trip_compartments.freight_total), 0)").
		Scan(&previousFreight)

	h.DB.Model(&models.TripCompartment{}).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trips.company_id = ? AND trips.created_at >= ? AND trips.created_at < ?", companyID, start, end).
		Select("COALESCE(SUM(trip_compartments.capacidade_litros), 0)").
		Scan(&totalLiters)

	h.DB.Model(&models.Trip{}).
		Where("company_id = ? AND status = ? AND created_at >= ? AND created_at < ?", companyID, models.StatusFinalizada, start, end).
		Count(&finishedTrips)

	h.DB.Model(&models.Trip{}).
		Where("company_id = ? AND status = ? AND created_at >= ? AND created_at < ?", companyID, models.StatusEmAndamento, start, end).
		Count(&activeTrips)

	var average float64
	if finishedTrips > 0 {
		average = totalFreight / float64(finishedTrips)
	}

	var variation float64
	if previousFreight > 0 {
		variation = ((totalFreight - previousFreight) / previousFreight) * 100
	}

	var byClient []ClientFreightReport

	h.DB.Model(&models.TripCompartment{}).
		Select(`
			trip_compartments.client_name_snapshot as client_name,
			COALESCE(SUM(trip_compartments.freight_total), 0) as total,
			COALESCE(SUM(trip_compartments.capacidade_litros), 0) as liters,
			COUNT(DISTINCT trip_compartments.trip_id) as trips
		`).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trips.company_id = ? AND trips.created_at >= ? AND trips.created_at < ?", companyID, start, end).
		Group("trip_compartments.client_name_snapshot").
		Order("total desc").
		Scan(&byClient)

	return ReportSummary{
		PeriodLabel:      start.Format("02/01/2006") + " até " + end.AddDate(0, 0, -1).Format("02/01/2006"),
		TotalFreight:     totalFreight,
		PreviousFreight:  previousFreight,
		FreightVariation: variation,
		TotalLiters:      totalLiters,
		FinishedTrips:    finishedTrips,
		ActiveTrips:      activeTrips,
		AveragePerTrip:   average,
		FreightByClient:  byClient,
	}
}
