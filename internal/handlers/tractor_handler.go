package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type TractorHandler struct {
	DB *gorm.DB
}

func NewTractorHandler(db *gorm.DB) *TractorHandler {
	return &TractorHandler{DB: db}
}

func (h *TractorHandler) List(c *gin.Context) {
	var tractors []models.Tractor

	h.DB.Where("company_id = ? AND active = true", CurrentCompanyID(c)).
		Order("created_at desc").
		Find(&tractors)

	RenderPage(c, "tractors/list.html", gin.H{
		"Title":    "Cavalos",
		"Tractors": tractors,
	})
}

func (h *TractorHandler) NewForm(c *gin.Context) {
	RenderPage(c, "tractors/form.html", gin.H{
		"Title": "Novo cavalo",
	})
}

func (h *TractorHandler) Create(c *gin.Context) {
	tractor := models.Tractor{
		CompanyID: CurrentCompanyID(c),
		Model:     c.PostForm("model"),
		Plate:     c.PostForm("plate"),
		Owner:     c.PostForm("owner"),
		Status:    c.PostForm("status"),
		Notes:     c.PostForm("notes"),
		Active:    true,
	}

	if tractor.Status == "" {
		tractor.Status = "ativo"
	}

	if err := h.DB.Create(&tractor).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao cadastrar cavalo: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/tractors")
}

func (h *TractorHandler) EditForm(c *gin.Context) {
	var tractor models.Tractor

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&tractor, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "cavalo não encontrado")
		return
	}

	RenderPage(c, "tractors/form.html", gin.H{
		"Title":   "Editar cavalo",
		"Tractor": tractor,
	})
}

func (h *TractorHandler) Update(c *gin.Context) {
	var tractor models.Tractor

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&tractor, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "cavalo não encontrado")
		return
	}

	tractor.Model = c.PostForm("model")
	tractor.Plate = c.PostForm("plate")
	tractor.Owner = c.PostForm("owner")
	tractor.Status = c.PostForm("status")
	tractor.Notes = c.PostForm("notes")

	if tractor.Status == "" {
		tractor.Status = "ativo"
	}

	if err := h.DB.Save(&tractor).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao atualizar cavalo: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/tractors")
}

func (h *TractorHandler) Delete(c *gin.Context) {
	if err := h.DB.Model(&models.Tractor{}).
		Where("company_id = ? AND id = ?", CurrentCompanyID(c), c.Param("id")).
		Updates(map[string]interface{}{
			"active": false,
			"status": "inativo",
		}).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao desativar cavalo: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/tractors")
}
