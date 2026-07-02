package handlers

import (
	"strings"

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

	search := strings.TrimSpace(c.Query("q"))

	query := h.DB.
		Preload("Truck").
		Preload("Compartments").
		Where("status = ?", models.StatusFinalizada)

	if search != "" {
		like := "%" + search + "%"

		query = query.
			Joins("LEFT JOIN trucks ON trucks.id = trips.truck_id").
			Joins("LEFT JOIN trip_compartments ON trip_compartments.trip_id = trips.id").
			Where(`
				trucks.cavalo_placa ILIKE ?
				OR trucks.tanque_placa ILIKE ?
				OR trucks.cavalo_modelo ILIKE ?
				OR trip_compartments.cliente ILIKE ?
			`, like, like, like, like).
			Group("trips.id")
	}

	query.Order("finished_at desc").Find(&trips)

	RenderPage(c, "history/index.html", gin.H{
		"Title":  "Histórico",
		"Trips":  trips,
		"Search": search,
	})
}
