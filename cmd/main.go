package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"transportadora/internal/database"
	"transportadora/internal/routes"
	"transportadora/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.Connect()
	services.EnsureAdminUser(db)

	r := gin.Default()

	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "troque-essa-chave-em-producao"
	}

	store := cookie.NewStore([]byte(secret))
	r.Use(sessions.Sessions("transportadora_session", store))

	routes.Register(r, db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("servidor rodando em http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("erro ao iniciar o servidor: %v", err)
	}
}
