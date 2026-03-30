package services

import (
	"context"
	"errors"
	"fmt"
	"turnos-medicos/internal/features/medicos/dto"
	"turnos-medicos/internal/pkg"

	medicoModel "turnos-medicos/internal/features/medicos/models"
	userModel "turnos-medicos/internal/features/users/models"

	medicoRepo "turnos-medicos/internal/features/medicos/repository"
	userRepo "turnos-medicos/internal/features/users/repository"
)

type MedicoService interface {
	CrearMedico(ctx context.Context, req dto.CrearMedicoRequest) (*medicoModel.Medico, error)
	ObtenerMedicoPorID(ctx context.Context, medicoID int64) (*medicoModel.Medico, error)
	ObtenerMedicoPorMatricula(ctx context.Context, matricula string) (*medicoModel.Medico, error)
	ListarMedicosActivos(ctx context.Context) ([]*medicoModel.Medico, error)
	ListarMedicosPorEspecialidad(ctx context.Context, especialidad string) ([]*medicoModel.Medico, error)
	ActualizarMedico(ctx context.Context, medicoID int64, req dto.ActualizarMedicoRequest) (*medicoModel.Medico, error)
	DesactivarMedico(ctx context.Context, medicoID int64) error
}

type medicoService struct {
	medicoRepo medicoRepo.MedicoRepository
	userRepo   userRepo.UserRepository
}

func NewMedicoService(medicoRepo medicoRepo.MedicoRepository, userRepo userRepo.UserRepository) MedicoService {
	return &medicoService{
		medicoRepo: medicoRepo,
		userRepo:   userRepo,
	}
}

// LÓGICA DE NEGOCIO

func (s *medicoService) CrearMedico(ctx context.Context, req dto.CrearMedicoRequest) (*medicoModel.Medico, error) {

	existeEmail, err := s.userRepo.ObtenerUsuarioPorEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
		return nil, err
	}
	if existeEmail != nil {
		return nil, pkg.ErrEmailDuplicado
	}

	existeMatricula, err := s.medicoRepo.ObtenerMedicoPorMatricula(ctx, req.Matricula)
	if err != nil && !errors.Is(err, pkg.ErrMedicoNoEncontrado) {
		return nil, err
	}
	if existeMatricula != nil {
		return nil, pkg.ErrMatriculaDuplicada
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	//construir entidad
	user := &userModel.User{
		Nombre:   req.Nombre,
		Apellido: req.Apellido,
		Email:    req.Email,
		Password: hashedPassword,
		Rol:      userModel.RolMedico,
		Activo:   true,
	}

	medico := &medicoModel.Medico{
		Especialidad: req.Especialidad,
		Matricula:    req.Matricula,
	}

	if err := s.medicoRepo.CrearMedico(ctx, user, medico); err != nil {
		return nil, err
	}

	medico.Nombre = user.Nombre
	medico.Apellido = user.Apellido
	medico.Email = user.Email
	medico.Activo = user.Activo

	return medico, nil

}

func (s *medicoService) ObtenerMedicoPorID(ctx context.Context, medicoID int64) (*medicoModel.Medico, error) {
	if medicoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	medico, err := s.medicoRepo.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if !medico.Activo {
		return nil, pkg.ErrMedicoInactivo
	}

	return medico, err
}

func (s *medicoService) ObtenerMedicoPorMatricula(ctx context.Context, matricula string) (*medicoModel.Medico, error) {

	if matricula == "" {
		return nil, pkg.ErrMatriculaRequerida
	}

	medico, err := s.medicoRepo.ObtenerMedicoPorMatricula(ctx, matricula)
	if err != nil {
		return nil, err
	}

	return medico, nil
}

func (s *medicoService) ListarMedicosActivos(ctx context.Context) ([]*medicoModel.Medico, error) {
	medicos, err := s.medicoRepo.ListarMedicosActivos(ctx)
	if err != nil {
		return nil, fmt.Errorf("error en servicio al listar medicos activos: %w", err)
	}
	return medicos, nil
}

func (s *medicoService) ListarMedicosPorEspecialidad(ctx context.Context, especialidad string) ([]*medicoModel.Medico, error) {

	if especialidad == "" {
		return nil, pkg.ErrEspecialidadRequerida
	}

	medicos, err := s.medicoRepo.ListarMedicosPorEspecialidad(ctx, especialidad)
	if err != nil {
		return nil, err
	}

	return medicos, nil

}

func (s *medicoService) ActualizarMedico(ctx context.Context, medicoID int64, req dto.ActualizarMedicoRequest) (*medicoModel.Medico, error) {
	if medicoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	medicoActual, err := s.medicoRepo.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if !medicoActual.Activo {
		return nil, pkg.ErrMedicoInactivo
	}

	// solo actualizar campos que vienen en el request
	if req.Nombre != "" {
		medicoActual.Nombre = req.Nombre
	}

	if req.Apellido != "" {
		medicoActual.Apellido = req.Apellido
	}

	if req.Email != "" && req.Email != medicoActual.Email {
		existeEmail, err := s.userRepo.ObtenerUsuarioPorEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
			return nil, err
		}
		if existeEmail != nil {
			return nil, pkg.ErrEmailDuplicado
		}
		medicoActual.Email = req.Email
	}
	if err := s.medicoRepo.ActualizarMedico(ctx, medicoActual); err != nil {
		return nil, err
	}

	return medicoActual, nil
}

func (s *medicoService) DesactivarMedico(ctx context.Context, medicoID int64) error {

	if err := s.medicoRepo.DesactivarMedico(ctx, medicoID); err != nil {
		return fmt.Errorf("error desactivando medico: %w", err)
	}
	return nil
}
