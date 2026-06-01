package main

import (
	"log"
	"turnos-medicos/internal/server"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Println(".env no encontrado, usando variables de entorno del sistema")
	}
	server.Start()

}
