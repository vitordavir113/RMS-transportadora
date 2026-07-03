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
	companyID := CurrentCompanyID(c)

	var tractors []models.Tractor
	var trailers []models.Trailer
	var drivers []models.Driver
	var clients []models.Client

	h.DB.Where("company_id = ? AND active = true", companyID).
		Order("plate asc").
		Find(&tractors)

	h.DB.Preload("Compartments", func(db *gorm.DB) *gorm.DB {
		return db.Order("numero asc")
	}).
		Where("company_id = ? AND active = true", companyID).
		Order("plate asc").
		Find(&trailers)

	h.DB.Where("company_id = ? AND active = true", companyID).
		Order("name asc").
		Find(&drivers)

	h.DB.Where("company_id = ? AND active = true", companyID).
		Order("name asc").
		Find(&clients)

	RenderPage(c, "trips/new.html", gin.H{
		"Title":    "Nova Viagem",
		"Tractors": tractors,
		"Trailers": trailers,
		"Drivers":  drivers,
		"Clients":  clients,
	})
}

func (h *TripHandler) LoadCompartments(c *gin.Context) {
	trailerID := c.Query("trailer_id")
	if trailerID == "" {
		c.String(http.StatusOK, "")
		return
	}

	var trailer models.Trailer
	if err := h.DB.Preload("Compartments", func(db *gorm.DB) *gorm.DB {
		return db.Order("numero asc")
	}).
		Where("company_id = ? AND active = true", CurrentCompanyID(c)).
		First(&trailer, trailerID).Error; err != nil {
		c.String(http.StatusNotFound, "tanque/carreta não encontrado")
		return
	}

	var clients []models.Client
	h.DB.Where("company_id = ? AND active = true", CurrentCompanyID(c)).
		Order("name asc").
		Find(&clients)

	RenderFragment(c, "trips/bocas_fragment.html", "bocas_fragment", gin.H{
		"Trailer":  trailer,
		"Clients":  clients,
		"Produtos": models.ProdutosDisponiveis,
	})
}

func (h *TripHandler) Create(c *gin.Context) {
	companyID := CurrentCompanyID(c)

	tractorID, err := strconv.Atoi(c.PostForm("tractor_id"))
	if err != nil || tractorID == 0 {
		c.String(http.StatusBadRequest, "cavalo inválido")
		return
	}

	trailerID, err := strconv.Atoi(c.PostForm("trailer_id"))
	if err != nil || trailerID == 0 {
		c.String(http.StatusBadRequest, "tanque/carreta inválido")
		return
	}

	driverID, err := strconv.Atoi(c.PostForm("driver_id"))
	if err != nil || driverID == 0 {
		c.String(http.StatusBadRequest, "motorista inválido")
		return
	}

	compIDs := c.PostFormArray("compartment_id[]")
	clientIDs := c.PostFormArray("client_id[]")
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

		clientID := 0
		produto := ""

		if i < len(clientIDs) {
			clientID, _ = strconv.Atoi(clientIDs[i])
		}

		if i < len(produtos) {
			produto = produtos[i]
		}

		if clientID == 0 || produto == "" {
			c.String(http.StatusBadRequest, "preencha cliente e produto em todas as bocas")
			return
		}

		bocas = append(bocas, services.BocaInput{
			TrailerCompartmentID: uint(cid),
			ClientID:             uint(clientID),
			Produto:              produto,
		})
	}

	trip, err := h.Service.CreateTrip(companyID, uint(tractorID), uint(trailerID), uint(driverID), bocas)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, "erro ao criar viagem: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trips/"+strconv.Itoa(int(trip.ID)))
}

func (h *TripHandler) Show(c *gin.Context) {
	id := c.Param("id")

	var trip models.Trip
	if err := h.DB.
		Preload("Compartments", func(db *gorm.DB) *gorm.DB {
			return db.Order("numero asc")
		}).
		Where("company_id = ?", CurrentCompanyID(c)).
		First(&trip, id).Error; err != nil {

		c.String(http.StatusNotFound, "viagem não encontrada")
		return
	}

	totalFrete := 0.0
	totalLitros := 0.0

	for _, c := range trip.Compartments {
		totalFrete += c.FreightTotal
		totalLitros += c.CapacidadeLitros
	}

	RenderPage(c, "trips/show.html", gin.H{
		"Title":       "Viagem #" + id,
		"Trip":        trip,
		"TotalFrete":  totalFrete,
		"TotalLitros": totalLitros,
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
		c.String(http.StatusInternalServerError, "erro ao listar viagens: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/history")

}

func (h *TripHandler) List(c *gin.Context) {
	var trips []models.Trip

	if err := h.DB.
		Preload("Compartments").
		Where("company_id = ?", CurrentCompanyID(c)).
		Order("created_at desc").
		Find(&trips).Error; err != nil {

		c.String(http.StatusInternalServerError, "erro ao listar viagens")
		return
	}

	type TripView struct {
		models.Trip
		TotalFrete  float64
		TotalLitros float64
	}

	var viewTrips []TripView

	for _, t := range trips {

		v := TripView{
			Trip: t,
		}

		for _, b := range t.Compartments {
			v.TotalFrete += b.FreightTotal
			v.TotalLitros += b.CapacidadeLitros
		}

		viewTrips = append(viewTrips, v)
	}

	RenderPage(c, "trips/list.html", gin.H{
		"Title": "Viagens",
		"Trips": viewTrips,
	})
}
