package bootstrap

import (
	"context"
	"log"
	"os"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/features/users/repository"
	"turnos-medicos/internal/pkg"
)

func SeedAdminUser(ctx context.Context, repo repository.UserRepository) {

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("ADMIN_EMAIL o ADMIN_PASSWORD no definidos. Seed cancelado.")
		return
	}

	// Verificar si el admin ya existe
	_, err := repo.ObtenerUsuarioPorEmail(ctx, email)
	if err == nil {
		log.Println("Usuario admin ya existe. Seed omitido.")
		return
	}

	// 2. Hashear password
	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		log.Println("Error hasheando password:", err)
		return
	}

	// 3. Crear usuario admin
	admin := &models.User{
		Nombre:   "Admin",
		Apellido: "Sistema",
		Email:    email,
		Password: hashedPassword,
		Rol:      models.RolAdmin,
		Activo:   true,
	}

	//Guardamos
	err = repo.CrearUsuario(ctx, admin)
	if err != nil {
		log.Println("Error creando usuario admin:", err)
		return
	}

	log.Println("Usuario admin creado correctamente.")

}
