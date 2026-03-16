package server

import (
	"fmt"
	"log"
	"os"
	"time"
	"turnos-medicos/internal/bootstrap"
	"turnos-medicos/internal/database"
	"turnos-medicos/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {

	database.Connect()
	defer database.DB.Close()

	handlers := bootstrap.Bootstrap(database.DB)

	r := gin.New()

	r.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			MaxAge:       12 * time.Hour,
		}),
	)

	routes.SetupRoutes(r, handlers)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)

	log.Fatal(r.Run(":" + port))
}
