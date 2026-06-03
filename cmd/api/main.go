package main

import (
	"log"
	"turnos-medicos/internal/server"

	"github.com/joho/godotenv"
)

// @title Turnos Médicos API
// @version 1.0
// @description API REST para gestión de turnos médicos.
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	_ = godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Println(".env no encontrado, usando variables de entorno del sistema")
	}
	server.Start()
}
