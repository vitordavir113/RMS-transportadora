package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type TruckHandler struct {
	DB *gorm.DB
}

func NewTruckHandler(db *gorm.DB) *TruckHandler {
	return &TruckHandler{DB: db}
}

func (h *TruckHandler) List(c *gin.Context) {
	var trucks []models.Truck

	h.DB.Preload("Compartments").
		Where("company_id = ?", CurrentCompanyID(c)).
		Order("created_at desc").
		Find(&trucks)

	RenderPage(c, "trucks/list.html", gin.H{
		"Title":  "Caminhões",
		"Trucks": trucks,
	})
}

func (h *TruckHandler) NewForm(c *gin.Context) {
	RenderPage(c, "trucks/form.html", gin.H{
		"Title": "Novo caminhão",
	})
}

func (h *TruckHandler) Create(c *gin.Context) {
	truck := models.Truck{
		CompanyID:    CurrentCompanyID(c),
		CavaloModelo: c.PostForm("cavalo_modelo"),
		CavaloPlaca:  c.PostForm("cavalo_placa"),
		TanquePlaca:  c.PostForm("tanque_placa"),
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&truck).Error; err != nil {
			return err
		}

		numeros := c.PostFormArray("boca_numero[]")
		capacidades := c.PostFormArray("boca_capacidade[]")

		for i := range numeros {
			if numeros[i] == "" || capacidades[i] == "" {
				continue
			}

			numero, err := strconv.Atoi(numeros[i])
			if err != nil {
				return err
			}

			capacidade, err := strconv.ParseFloat(capacidades[i], 64)
			if err != nil {
				return err
			}

			compartment := models.TruckCompartment{
				TruckID:          truck.ID,
				Numero:           numero,
				CapacidadeLitros: capacidade,
			}

			if err := tx.Create(&compartment).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "erro ao cadastrar caminhão: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trucks")
}

func (h *TruckHandler) EditForm(c *gin.Context) {
	id := c.Param("id")

	var truck models.Truck
	if err := h.DB.Preload("Compartments").
		Where("company_id = ?", CurrentCompanyID(c)).
		First(&truck, id).Error; err != nil {
		c.String(http.StatusNotFound, "caminhão não encontrado")
		return
	}

	RenderPage(c, "trucks/edit.html", gin.H{
		"Title": "Editar caminhão",
		"Truck": truck,
	})
}

func (h *TruckHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var truck models.Truck
	if err := h.DB.
		Where("company_id = ?", CurrentCompanyID(c)).
		First(&truck, id).Error; err != nil {
		c.String(http.StatusNotFound, "caminhão não encontrado")
		return
	}

	truck.CavaloModelo = c.PostForm("cavalo_modelo")
	truck.CavaloPlaca = c.PostForm("cavalo_placa")
	truck.TanquePlaca = c.PostForm("tanque_placa")

	if err := h.DB.Save(&truck).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao atualizar caminhão: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trucks/"+id+"/edit")
}

func (h *TruckHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	companyID := CurrentCompanyID(c)

	var count int64
	h.DB.Model(&models.Trip{}).
		Where("truck_id = ? AND company_id = ?", id, companyID).
		Count(&count)

	if count > 0 {
		c.String(http.StatusBadRequest, "não é possível excluir caminhão com viagens cadastradas")
		return
	}

	if err := h.DB.
		Where("company_id = ?", companyID).
		Delete(&models.Truck{}, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao excluir caminhão: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trucks")
}

func (h *TruckHandler) AddCompartment(c *gin.Context) {
	id := c.Param("id")
	companyID := CurrentCompanyID(c)

	var truck models.Truck
	if err := h.DB.
		Where("company_id = ?", companyID).
		First(&truck, id).Error; err != nil {
		c.String(http.StatusNotFound, "caminhão não encontrado")
		return
	}

	numero, err := strconv.Atoi(c.PostForm("numero"))
	if err != nil {
		c.String(http.StatusBadRequest, "número da boca inválido")
		return
	}

	capacidade, err := strconv.ParseFloat(c.PostForm("capacidade_litros"), 64)
	if err != nil {
		c.String(http.StatusBadRequest, "capacidade inválida")
		return
	}

	compartment := models.TruckCompartment{
		TruckID:          truck.ID,
		Numero:           numero,
		CapacidadeLitros: capacidade,
	}

	if err := h.DB.Create(&compartment).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao adicionar boca: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trucks/"+id+"/edit")
}

func (h *TruckHandler) DeleteCompartment(c *gin.Context) {
	truckID := c.Param("id")
	compID := c.Param("compId")
	companyID := CurrentCompanyID(c)

	var truck models.Truck
	if err := h.DB.
		Where("id = ? AND company_id = ?", truckID, companyID).
		First(&truck).Error; err != nil {
		c.String(http.StatusNotFound, "caminhão não encontrado")
		return
	}

	var count int64
	h.DB.Model(&models.TripCompartment{}).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trip_compartments.truck_compartment_id = ? AND trips.company_id = ?", compID, companyID).
		Count(&count)

	if count > 0 {
		c.String(http.StatusBadRequest, "não é possível excluir boca usada em viagem")
		return
	}

	if err := h.DB.
		Where("truck_id = ?", truck.ID).
		Delete(&models.TruckCompartment{}, compID).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao excluir boca: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trucks/"+truckID+"/edit")
}

func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}
