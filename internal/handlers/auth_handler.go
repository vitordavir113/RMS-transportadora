package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) LoginForm(c *gin.Context) {
	RenderAuthPage(c, "auth/login.html", gin.H{
		"Title": "Login",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := h.DB.Where("username = ?", username).First(&user).Error; err != nil {
		RenderAuthPage(c, "auth/login.html", gin.H{
			"Title": "Login",
			"Error": "Usuário ou senha inválidos",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		RenderAuthPage(c, "auth/login.html", gin.H{
			"Title": "Login",
			"Error": "Usuário ou senha inválidos",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("company_id", user.CompanyID)
	session.Set("username", user.Username)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusSeeOther, "/login")
}
