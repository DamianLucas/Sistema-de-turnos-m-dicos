package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {

	//leer variables de entorno
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//Validar que existan
	if host == "" || portStr == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ Faltan variables de entorno de la base de datos")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("❌ DB_PORT debe ser un número")
	}

	//Crear string de conexion
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Error al hacer ping a la base de datos:", err)
	}

	fmt.Println("✅ Conectado a PostgreSQL")
}
