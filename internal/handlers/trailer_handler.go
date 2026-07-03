package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type TrailerHandler struct {
	DB *gorm.DB
}

func NewTrailerHandler(db *gorm.DB) *TrailerHandler {
	return &TrailerHandler{DB: db}
}

func (h *TrailerHandler) List(c *gin.Context) {
	var trailers []models.Trailer

	h.DB.Preload("Compartments").
		Where("company_id = ? AND active = true", CurrentCompanyID(c)).
		Order("created_at desc").
		Find(&trailers)

	RenderPage(c, "trailers/list.html", gin.H{
		"Title":    "Tanques e Carretas",
		"Trailers": trailers,
	})
}

func (h *TrailerHandler) NewForm(c *gin.Context) {
	RenderPage(c, "trailers/form.html", gin.H{
		"Title": "Novo tanque/carreta",
	})
}

func (h *TrailerHandler) Create(c *gin.Context) {
	trailer := models.Trailer{
		CompanyID: CurrentCompanyID(c),
		Plate:     c.PostForm("plate"),
		Type:      c.PostForm("type"),
		Owner:     c.PostForm("owner"),
		Status:    c.PostForm("status"),
		Notes:     c.PostForm("notes"),
		Active:    true,
	}

	if trailer.Type == "" {
		trailer.Type = "tanque"
	}

	if trailer.Status == "" {
		trailer.Status = "ativo"
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trailer).Error; err != nil {
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

			compartment := models.TrailerCompartment{
				TrailerID:        trailer.ID,
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
		c.String(http.StatusInternalServerError, "erro ao cadastrar tanque/carreta: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trailers")
}

func (h *TrailerHandler) EditForm(c *gin.Context) {
	var trailer models.Trailer

	if err := h.DB.Preload("Compartments").
		Where("company_id = ?", CurrentCompanyID(c)).
		First(&trailer, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "tanque/carreta não encontrado")
		return
	}

	RenderPage(c, "trailers/edit.html", gin.H{
		"Title":   "Editar tanque/carreta",
		"Trailer": trailer,
	})
}

func (h *TrailerHandler) Update(c *gin.Context) {
	var trailer models.Trailer

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&trailer, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "tanque/carreta não encontrado")
		return
	}

	trailer.Plate = c.PostForm("plate")
	trailer.Type = c.PostForm("type")
	trailer.Owner = c.PostForm("owner")
	trailer.Status = c.PostForm("status")
	trailer.Notes = c.PostForm("notes")

	if trailer.Type == "" {
		trailer.Type = "tanque"
	}

	if trailer.Status == "" {
		trailer.Status = "ativo"
	}

	if err := h.DB.Save(&trailer).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao atualizar tanque/carreta: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trailers")
}

func (h *TrailerHandler) Delete(c *gin.Context) {
	if err := h.DB.Model(&models.Trailer{}).
		Where("company_id = ? AND id = ?", CurrentCompanyID(c), c.Param("id")).
		Updates(map[string]interface{}{
			"active": false,
			"status": "inativo",
		}).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao desativar tanque/carreta: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trailers")
}

func (h *TrailerHandler) AddCompartment(c *gin.Context) {
	var trailer models.Trailer

	if err := h.DB.Where("id = ? AND company_id = ?", c.Param("id"), CurrentCompanyID(c)).
		First(&trailer).Error; err != nil {
		c.String(http.StatusNotFound, "tanque/carreta não encontrado")
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

	compartment := models.TrailerCompartment{
		TrailerID:        trailer.ID,
		Numero:           numero,
		CapacidadeLitros: capacidade,
	}

	if err := h.DB.Create(&compartment).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao adicionar boca: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trailers/"+c.Param("id")+"/edit")
}

func (h *TrailerHandler) DeleteCompartment(c *gin.Context) {
	trailerID := c.Param("id")
	compID := c.Param("compId")
	companyID := CurrentCompanyID(c)

	var count int64
	h.DB.Model(&models.TripCompartment{}).
		Joins("JOIN trips ON trips.id = trip_compartments.trip_id").
		Where("trip_compartments.trailer_compartment_id = ? AND trips.company_id = ?", compID, companyID).
		Count(&count)

	if count > 0 {
		c.String(http.StatusBadRequest, "não é possível excluir boca usada em viagem")
		return
	}

	if err := h.DB.
		Where("trailer_id = ?", trailerID).
		Delete(&models.TrailerCompartment{}, compID).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao excluir boca: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/trailers/"+trailerID+"/edit")
}
