package services

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"transportadora/internal/models"
)

func EnsureAdminUser(db *gorm.DB) {
	createCompanyWithUser(db, "Transportadora A", "voltadajurema", "Oleo24250")
	createCompanyWithUser(db, "Transportadora B", "vitor", "vvddmmrr")
}

func createCompanyWithUser(db *gorm.DB, companyName, username, password string) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err == nil {
		return
	}

	company := models.Company{Name: companyName}

	if err := db.Create(&company).Error; err != nil {
		log.Fatalf("erro ao criar empresa %s: %v", companyName, err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("erro ao gerar senha: %v", err)
	}

	user = models.User{
		CompanyID:    company.ID,
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("erro ao criar usuário %s: %v", username, err)
	}

	log.Printf("usuário criado: %s / %s | empresa: %s", username, password, companyName)
}
