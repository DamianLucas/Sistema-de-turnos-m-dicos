package bootstrap

import (
	"context"
	"database/sql"
	handlerAuth "turnos-medicos/internal/features/auth/handlers"
	handlerUser "turnos-medicos/internal/features/users/handlers"

	newUserPostgresRepo "turnos-medicos/internal/features/users/repository/postgres"
	newUserService "turnos-medicos/internal/features/users/services"

	newAuthService "turnos-medicos/internal/features/auth/service"
)

type Handlers struct {
	User *handlerUser.UserHandler
	Auth *handlerAuth.AuthHandler
}

func Bootstrap(db *sql.DB) *Handlers {
	//Repositories
	userRepo := newUserPostgresRepo.NewUserPostgresRepository(db)

	//services
	userService := newUserService.NewUserService(userRepo)
	authService := newAuthService.NewAuthService(userRepo)

	//handlers
	userHandler := handlerUser.NewUserHandler(userService)
	authHandler := handlerAuth.NewAuthHandler(authService)

	//Seed
	SeedAdminUser(context.Background(), userService)

	return &Handlers{
		User: userHandler,
		Auth: authHandler,
	}

}
