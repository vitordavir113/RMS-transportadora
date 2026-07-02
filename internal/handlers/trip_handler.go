package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
	"transportadora/internal/services"
)

type TripHandler struct {
	DB      *gorm.DB
	Service *services.TripService
}

func NewTripHandler(db *gorm.DB) *TripHandler {
	return &TripHandler{
		DB:      db,
		Service: services.NewTripService(db),
	}
}

func (h *TripHandler) NewForm(c *gin.Context) {
	var trucks []models.Truck
	companyID := CurrentCompanyID(c)

	h.DB.Preload("Compartments", func(db *gorm.DB) *gorm.DB {
		return db.Order("numero asc")
	}).
		Where("company_id = ?", companyID).
		Order("cavalo_modelo asc").
		Find(&trucks)

	RenderPage(c, "trips/new.html", gin.H{
		"Title":  "Nova Viagem",
		"Trucks": trucks,
	})
}

func (h *TripHandler) LoadCompartments(c *gin.Context) {
	truckID := c.Query("truck_id")
	if truckID == "" {
		c.String(http.StatusOK, "")
		return
	}

	var truck models.Truck
	if err := h.DB.Preload("Compartments", func(db *gorm.DB) *gorm.DB {
		return db.Order("numero asc")
	}).
		Where("company_id = ?", CurrentCompanyID(c)).
		First(&truck, truckID).Error; err != nil {
		c.String(http.StatusNotFound, "caminhão não encontrado")
		return
	}

	RenderFragment(c, "trips/bocas_fragment.html", "bocas_fragment", gin.H{
		"Truck":    truck,
		"Produtos": models.ProdutosDisponiveis,
	})
}

func (h *TripHandler) Create(c *gin.Context) {
	companyID := CurrentCompanyID(c)

	truckID, err := strconv.Atoi(c.PostForm("truck_id"))
	if err != nil || truckID == 0 {
		c.String(http.StatusBadRequest, "caminhão inválido")
		return
	}

	compIDs := c.PostFormArray("compartment_id[]")
	clientes := c.PostFormArray("cliente[]")
	produtos := c.PostFormArray("produto[]")

	if len(compIDs) == 0 {
		c.String(http.StatusBadRequest, "nenhuma boca foi carregada")
		return
	}

	var bocas []services.BocaInput

	for i := range compIDs {
		cid, err := strconv.Atoi(compIDs[i])
		if err != nil || cid == 0 {
			continue
		}

		cliente := ""
		produto := ""

		if i < len(clientes) {
			cliente = clientes[i]
		}

		if i < len(produtos) {
			produto = produtos[i]
		}

		if cliente == "" || produto == "" {
			c.String(http.StatusBadRequest, "preencha cliente e produto em todas as bocas")
			return
		}

		bocas = append(bocas, services.BocaInput{
			TruckCompartmentID: uint(cid),
			Cliente:            cliente,
			Produto:            produto,
		})
	}

	trip, err := h.Service.CreateTrip(companyID, uint(truckID), bocas)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, "erro ao criar viagem: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trips/"+strconv.Itoa(int(trip.ID)))
}

func (h *TripHandler) Show(c *gin.Context) {
	id := c.Param("id")

	var trip models.Trip
	if err := h.DB.Preload("Truck").
		Preload("Compartments", func(db *gorm.DB) *gorm.DB {
			return db.Order("numero asc")
		}).
		Where("company_id = ?", CurrentCompanyID(c)).
		First(&trip, id).Error; err != nil {
		c.String(http.StatusNotFound, "viagem não encontrada")
		return
	}

	RenderPage(c, "trips/show.html", gin.H{
		"Title":              "Viagem #" + id,
		"Trip":               trip,
		"TotalLitros":        trip.TotalLitros(),
		"QuantidadeClientes": trip.QuantidadeClientes(),
	})
}

func (h *TripHandler) Finish(c *gin.Context) {
	id := c.Param("id")

	tripID, err := strconv.Atoi(id)
	if err != nil || tripID == 0 {
		c.String(http.StatusBadRequest, "viagem inválida")
		return
	}

	if err := h.Service.FinishTrip(CurrentCompanyID(c), uint(tripID)); err != nil {
		c.String(http.StatusUnprocessableEntity, "erro ao finalizar viagem: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/history")
}
