package main

import (
	"log"
	"turnos-medicos/internal/server"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error al cargar archivo .env")
	}

	server.Start()

}
