package services

import (
	"context"
	"fmt"

	"turnos-medicos/internal/features/users/dto"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/features/users/repository"
	"turnos-medicos/internal/utils"
)

type UserService interface {
	CrearUsuario(ctx context.Context, req dto.CrearUsuarioRequest) (*models.User, error)
	ObtenerUsuarioPorID(ctx context.Context, id int64) (*models.User, error)
	ListarUsuariosActivos(ctx context.Context) ([]*models.User, error)
	DesactivarUsuario(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// LÓGICA DE NEGOCIO

func (s *userService) CrearUsuario(ctx context.Context, req dto.CrearUsuarioRequest) (*models.User, error) {

	//validaciones de negocio

	//validar email unico
	existeEmail, err := s.repo.ObtenerUsuarioPorEmail(ctx, req.Email)
	if err == nil && existeEmail != nil {
		return nil, utils.ErrEmailDuplicado
	}

	//hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	//construir entidad
	user := &models.User{
		Nombre:   req.Nombre,
		Apellido: req.Apellido,
		Email:    req.Email,
		Password: hashedPassword,
		Rol:      req.Rol,
		Activo:   true,
	}
	if err := s.repo.CrearUsuario(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ObtenerUsuarioPorID(ctx context.Context, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, utils.ErrIDInvalido
	}

	user, err := s.repo.ObtenerUsuarioPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !user.Activo {
		return nil, utils.ErrUsuarioInactivo
	}

	return user, err

}

func (s *userService) ListarUsuariosActivos(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.ListarUsuariosActivos(ctx)
	if err != nil {
		return nil, fmt.Errorf("error en servicio al listar usuarios activos: %w", err)
	}
	return users, nil
}

func (s *userService) DesactivarUsuario(ctx context.Context, id int64) error {
	user, err := s.repo.ObtenerUsuarioPorID(ctx, id)
	if err != nil {
		return err
	}

	//si ya esta inactivo no se hace nada
	if !user.Activo {
		return nil
	}

	if err := s.repo.DesactivarUsuario(ctx, id); err != nil {
		return fmt.Errorf("error desactivando usuario: %w", err)
	}

	return nil

}
