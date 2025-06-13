package main

import (
	"log"

	_ "thothix-backend/docs" // Importa i documenti Swagger generati
	"thothix-backend/internal/config"
	"thothix-backend/internal/database"
	"thothix-backend/internal/router"
)

// @title Thothix API
// @version 1.0
// @description API per la piattaforma di messaggistica aziendale Thothix
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.thothix.com/support
// @contact.email support@thothix.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:30000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Carica configurazione
	cfg := config.Load()

	// Inizializza database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Esegui migrazioni
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Inizializza router
	r := router.Setup(db, cfg)

	// Avvia server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
