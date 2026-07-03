package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type HistoryHandler struct {
	DB *gorm.DB
}

func NewHistoryHandler(db *gorm.DB) *HistoryHandler {
	return &HistoryHandler{DB: db}
}

func (h *HistoryHandler) Index(c *gin.Context) {
	var trips []models.Trip

	if err := h.DB.
		Preload("Compartments").
		Where("company_id = ? AND status = ?", CurrentCompanyID(c), models.StatusFinalizada).
		Order("finished_at desc").
		Find(&trips).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao carregar histórico")
		return
	}

	type TripHistoryView struct {
		models.Trip
		TotalFrete  float64
		TotalLitros float64
	}

	var viewTrips []TripHistoryView

	for _, trip := range trips {
		item := TripHistoryView{
			Trip: trip,
		}

		for _, b := range trip.Compartments {
			item.TotalFrete += b.FreightTotal
			item.TotalLitros += b.CapacidadeLitros
		}

		viewTrips = append(viewTrips, item)
	}

	RenderPage(c, "history/index.html", gin.H{
		"Title": "Histórico",
		"Trips": viewTrips,
	})
}
