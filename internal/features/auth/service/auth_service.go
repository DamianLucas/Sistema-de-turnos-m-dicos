package service

import (
	"context"
	"errors"
	"turnos-medicos/internal/features/auth/dto"
	"turnos-medicos/internal/features/users/repository"
	"turnos-medicos/internal/utils"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*dto.LoginResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Login(ctx context.Context, email string, password string) (*dto.LoginResponse, error) {
	user, err := s.repo.ObtenerUsuarioPorEmail(ctx, email)
	if err != nil {
		if errors.Is(err, utils.ErrUsuarioNoEncontrado) {
			return nil, utils.ErrCredencialesInvalidas
		}
		return nil, err
	}

	if !user.Activo {
		return nil, utils.ErrUsuarioInactivo
	}

	if !utils.VerificarPassword(password, user.Password) {
		return nil, utils.ErrCredencialesInvalidas
	}

	token, err := utils.GenerarToken(user.ID, user.Rol)
	if err != nil {
		return nil, errors.New("error al generar token")
	}

	var resp dto.LoginResponse
	resp.Token = token
	resp.User.ID = int(user.ID)
	resp.User.Email = user.Email
	resp.User.Rol = string(user.Rol)

	return &resp, nil
}
