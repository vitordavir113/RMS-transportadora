package handlers

import (
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

func (h *HomeHandler) Index(c *gin.Context) {
	companyID := CurrentCompanyID(c)

	var trips []models.Trip

	h.DB.Preload("Truck").
		Preload("Compartments").
		Where("company_id = ? AND status = ?", companyID, models.StatusEmAndamento).
		Order("created_at desc").
		Find(&trips)

	var totalTrucks int64
	h.DB.Model(&models.Truck{}).
		Where("company_id = ?", companyID).
		Count(&totalTrucks)

	var totalLitros float64
	for _, trip := range trips {
		totalLitros += trip.TotalLitros()
	}

	RenderPage(c, "home/index.html", gin.H{
		"Title":       "Viagens em Andamento",
		"Trips":       trips,
		"TotalTrips":  len(trips),
		"TotalLitros": totalLitros,
		"TotalTrucks": totalTrucks,
	})
}
