package services

import (
	"context"
	"fmt"

	"turnos-medicos/internal/features/users/dto"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/features/users/repository"
	"turnos-medicos/internal/pkg"
)

type UserService interface {
	CrearUsuario(ctx context.Context, req dto.CrearUsuarioRequest) (*models.User, error)
	ObtenerUsuarioPorID(ctx context.Context, id int64) (*models.User, error)
	ListarUsuariosActivos(ctx context.Context) ([]*models.User, error)
	ActualizarUsuario(ctx context.Context, id int64, req dto.ActualizarUsuarioRequest) (*models.User, error)
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
		return nil, pkg.ErrEmailDuplicado
	}

	//hash password
	hashedPassword, err := pkg.HashPassword(req.Password)
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

func (s *userService) ObtenerUsuarioPorID(ctx context.Context, userID int64) (*models.User, error) {
	if userID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	user, err := s.repo.ObtenerUsuarioPorID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.Activo {
		return nil, pkg.ErrUsuarioInactivo
	}

	return user, nil

}

func (s *userService) ListarUsuariosActivos(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.ListarUsuariosActivos(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", pkg.ErrListarUsuariosActivos, err)
	}
	return users, nil
}

func (s *userService) ActualizarUsuario(ctx context.Context, id int64, req dto.ActualizarUsuarioRequest) (*models.User, error) {
	if id <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	userActual, err := s.repo.ObtenerUsuarioPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !userActual.Activo {
		return nil, pkg.ErrUsuarioInactivo
	}

	// solo actualizar campos que vienen en el request
	if req.Nombre != "" {
		userActual.Nombre = req.Nombre
	}
	if req.Apellido != "" {
		userActual.Apellido = req.Apellido
	}
	if req.Email != "" {
		userActual.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		userActual.Password = hashedPassword
	}
	if req.Rol != "" {
		userActual.Rol = req.Rol
	}

	if err := s.repo.ActualizarUsuario(ctx, userActual); err != nil {
		return nil, err
	}

	return userActual, nil
}

func (s *userService) DesactivarUsuario(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return pkg.ErrIDInvalido
	}

	user, err := s.repo.ObtenerUsuarioPorID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.Activo {
		return pkg.ErrUsuarioInactivo
	}

	if err := s.repo.DesactivarUsuario(ctx, userID); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrDesactivarUsuario, err)
	}

	return nil
}
