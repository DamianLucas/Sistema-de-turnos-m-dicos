package bootstrap

import (
	"context"
	"database/sql"

	//auth
	handlerAuth "turnos-medicos/internal/features/auth/handlers"
	newAuthService "turnos-medicos/internal/features/auth/service"

	//User
	handlerUser "turnos-medicos/internal/features/users/handlers"
	newUserPostgresRepo "turnos-medicos/internal/features/users/repository/postgres"
	newUserService "turnos-medicos/internal/features/users/services"

	//Medico
	handlerMedico "turnos-medicos/internal/features/medicos/handlers"
	newMedicoPostgresRepo "turnos-medicos/internal/features/medicos/repository/postgres"
	newMedicoService "turnos-medicos/internal/features/medicos/services"

	//Pacientes
	newPacientePostgresRepo "turnos-medicos/internal/features/pacientes/repository/postgres"
)

type Handlers struct {
	User   *handlerUser.UserHandler
	Auth   *handlerAuth.AuthHandler
	Medico *handlerMedico.MedicoHandler
}

func Bootstrap(db *sql.DB) *Handlers {

	//USERS
	//Repositories
	userRepo := newUserPostgresRepo.NewUserPostgresRepository(db)
	authService := newAuthService.NewAuthService(userRepo)

	//services
	userService := newUserService.NewUserService(userRepo)

	//handlers
	userHandler := handlerUser.NewUserHandler(userService)
	authHandler := handlerAuth.NewAuthHandler(authService)

	//MEDICOS
	//Repositories
	medicoRepo := newMedicoPostgresRepo.NewMedicoPostgresRepository(db)
	pacienteRepo := newPacientePostgresRepo.NewPacientePostgresRepository(db)

	//services
	medicoService := newMedicoService.NewMedicoService(medicoRepo, userRepo, pacienteRepo, db)

	//handlers
	medicoHandler := handlerMedico.NewMedicoHandler(medicoService)

	//Seed
	SeedAdminUser(context.Background(), userService)

	return &Handlers{
		User:   userHandler,
		Auth:   authHandler,
		Medico: medicoHandler,
	}

}
