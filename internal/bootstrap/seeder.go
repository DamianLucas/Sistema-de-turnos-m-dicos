package bootstrap

import (
	"context"
	"errors"
	"log"
	"os"

	"turnos-medicos/internal/features/users/dto"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/features/users/services"
	"turnos-medicos/internal/pkg"
)

func SeedAdminUser(ctx context.Context, userService services.UserService) {

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("ADMIN_EMAIL o ADMIN_PASSWORD no definidos. Seed cancelado.")
		return
	}

	createUserAdmin := dto.CrearUsuarioRequest{
		Nombre:   "Admin",
		Apellido: "Sistema",
		Email:    email,
		Password: password,
		Rol:      models.RolAdmin,
	}

	_, err := userService.CrearUsuario(ctx, createUserAdmin)
	if err != nil {

		if errors.Is(err, pkg.ErrEmailDuplicado) {
			log.Println("Usuario admin ya existe. Seed omitido.")
			return
		}

		log.Println("Error creando usuario admin:", err)
		return
	}

	log.Println("Usuario admin creado correctamente.")
}
