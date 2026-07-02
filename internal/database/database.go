package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"transportadora/internal/models"
)

var DB *gorm.DB

// Connect abre a conexão com o PostgreSQL usando variáveis de ambiente.
func Connect() *gorm.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "transportadora")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=America/Sao_Paulo",
		host, port, user, password, dbname, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("erro ao conectar no banco de dados: %v", err)
	}

	DB = db
	return db
}

// AutoMigrate cria/atualiza as tabelas automaticamente a partir dos models.
// As migrations SQL equivalentes ficam em /migrations para referência e
// para quem preferir aplicar manualmente via psql/migrate.
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Truck{},
		&models.TruckCompartment{},
		&models.Trip{},
		&models.TripCompartment{},
	)
	if err != nil {
		log.Fatalf("erro ao rodar as migrations: %v", err)
	}
	log.Println("migrations aplicadas com sucesso")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
