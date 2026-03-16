package repository

import (
	"context"
	"turnos-medicos/internal/features/users/models"
)

type UserRepository interface {
	CrearUsuario(ctx context.Context, u *models.User) error
	ListarUsuariosActivos(ctx context.Context) ([]*models.User, error)

	ObtenerUsuarioPorID(ctx context.Context, id int64) (*models.User, error)
	ObtenerUsuarioPorEmail(ctx context.Context, email string) (*models.User, error)
	ObtenerUsuarioPorRol(ctx context.Context, rol models.Rol) ([]*models.User, error)

	ActualizarUsuario(ctx context.Context, u *models.User) error

	DesactivarUsuario(ctx context.Context, id int64) error
}
