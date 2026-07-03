package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type DriverHandler struct {
	DB *gorm.DB
}

func NewDriverHandler(db *gorm.DB) *DriverHandler {
	return &DriverHandler{DB: db}
}

func (h *DriverHandler) List(c *gin.Context) {
	var drivers []models.Driver

	h.DB.Where("company_id = ? AND active = true", CurrentCompanyID(c)).
		Order("name asc").
		Find(&drivers)

	RenderPage(c, "drivers/list.html", gin.H{
		"Title":   "Motoristas",
		"Drivers": drivers,
	})
}

func (h *DriverHandler) NewForm(c *gin.Context) {
	RenderPage(c, "drivers/form.html", gin.H{
		"Title": "Novo motorista",
	})
}

func (h *DriverHandler) Create(c *gin.Context) {
	driver := models.Driver{
		CompanyID: CurrentCompanyID(c),
		Name:      c.PostForm("name"),
		CPF:       c.PostForm("cpf"),
		Phone:     c.PostForm("phone"),
		CNH:       c.PostForm("cnh"),
		Active:    true,
	}

	if err := h.DB.Create(&driver).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao cadastrar motorista: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/drivers")
}

func (h *DriverHandler) EditForm(c *gin.Context) {
	var driver models.Driver

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&driver, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "motorista não encontrado")
		return
	}

	RenderPage(c, "drivers/form.html", gin.H{
		"Title":  "Editar motorista",
		"Driver": driver,
	})
}

func (h *DriverHandler) Update(c *gin.Context) {
	var driver models.Driver

	if err := h.DB.Where("company_id = ?", CurrentCompanyID(c)).
		First(&driver, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "motorista não encontrado")
		return
	}

	driver.Name = c.PostForm("name")
	driver.CPF = c.PostForm("cpf")
	driver.Phone = c.PostForm("phone")
	driver.CNH = c.PostForm("cnh")

	if err := h.DB.Save(&driver).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao atualizar motorista: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/drivers")
}

func (h *DriverHandler) Delete(c *gin.Context) {
	if err := h.DB.Model(&models.Driver{}).
		Where("company_id = ? AND id = ?", CurrentCompanyID(c), c.Param("id")).
		Update("active", false).Error; err != nil {
		c.String(http.StatusInternalServerError, "erro ao desativar motorista: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/drivers")
}
