package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type ClientHandler struct {
	DB *gorm.DB
}

func NewClientHandler(db *gorm.DB) *ClientHandler {
	return &ClientHandler{DB: db}
}

func (h *ClientHandler) List(c *gin.Context) {
	var clients []models.Client

	h.DB.Where("company_id = ? AND active = true", CurrentCompanyID(c)).
		Order("name asc").
		Find(&clients)

	RenderPage(c, "clients/list.html", gin.H{
		"Title":   "Clientes",
		"Clients": clients,
	})
}

func (h *ClientHandler) NewForm(c *gin.Context) {
	RenderPage(c, "clients/form.html", gin.H{
		"Title": "Novo cliente",
	})
}

func (h *ClientHandler) Create(c *gin.Context) {
	client := models.Client{
		CompanyID:    CurrentCompanyID(c),
		Name:         c.PostForm("name"),
		Document:     c.PostForm("document"),
		Phone:        c.PostForm("phone"),
		FreightType:  models.FreightType(c.PostForm("freight_type")),
		FreightValue: parseFloat(c.PostForm("freight_value")),
		Active:       true,
	}

	if client.FreightType == "" {
		client.FreightType = models.FreightPorBoca
	}

	if err := h.DB.Create(&client).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao cadastrar cliente: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/clients")
}

func (h *ClientHandler) EditForm(c *gin.Context) {
	var client models.Client

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&client, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "cliente não encontrado")
		return
	}

	RenderPage(c, "clients/form.html", gin.H{
		"Title":  "Editar cliente",
		"Client": client,
	})
}

func (h *ClientHandler) Update(c *gin.Context) {
	var client models.Client

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&client, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "cliente não encontrado")
		return
	}

	client.Name = c.PostForm("name")
	client.Document = c.PostForm("document")
	client.Phone = c.PostForm("phone")
	client.FreightType = models.FreightType(c.PostForm("freight_type"))
	client.FreightValue = parseFloat(c.PostForm("freight_value"))

	if err := h.DB.Save(&client).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao atualizar cliente: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/clients")
}

func (h *ClientHandler) Delete(c *gin.Context) {
	if err := h.DB.Model(&models.Client{}).
		Where("company_id = ? AND id = ?", CurrentCompanyID(c), c.Param("id")).
		Update("active", false).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao desativar cliente: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/clients")
}
