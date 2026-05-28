package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"turnos-medicos/internal/features/pacientes/dto"
	"turnos-medicos/internal/features/pacientes/models"
	userModel "turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/pkg"

	repositoryMedico "turnos-medicos/internal/features/medicos/repository"
	repositoryPaciente "turnos-medicos/internal/features/pacientes/repository"
)

type PacienteService interface {
	CrearPaciente(ctx context.Context, req dto.CrearPacienteRequest) (*models.Paciente, error)
	ObtenerPacientePorID(ctx context.Context, authUserID int64, authRol userModel.Rol, pacienteID int64) (*models.Paciente, error)
	ObtenerPacientePorDNI(ctx context.Context, dni string) (*models.Paciente, error)
	ListarPacientesActivos(ctx context.Context, authUserID int64, authRol userModel.Rol) ([]*models.Paciente, error)
	DesactivarPaciente(ctx context.Context, pacienteID int64) error
	ActivarPaciente(ctx context.Context, pacienteID int64) error
	ActualizarPaciente(ctx context.Context, authUserID int64, authRol userModel.Rol, pacienteID int64, req dto.ActualizarPacienteRequest) (*models.Paciente, error)
	AsignarMedicoTratante(ctx context.Context, pacienteID, medicoID int64) error
	QuitarMedicoTratante(ctx context.Context, pacienteID int64) error
	ListarPacientesPorMedico(ctx context.Context, authUserID int64, authRol userModel.Rol, medicoID int64) ([]*models.Paciente, error)
}

type pacienteService struct {
	repoPaciente repositoryPaciente.PacienteRepository
	repoMedico   repositoryMedico.MedicoRepository
}

func NewPacienteService(repoPaciente repositoryPaciente.PacienteRepository, repoMedico repositoryMedico.MedicoRepository) PacienteService {
	return &pacienteService{
		repoPaciente: repoPaciente,
		repoMedico:   repoMedico,
	}
}

// LÓGICA DE NEGOCIO

func (s *pacienteService) CrearPaciente(ctx context.Context, req dto.CrearPacienteRequest) (*models.Paciente, error) {
	existe, err := s.repoPaciente.ObtenerPacientePorDNI(ctx, req.DNI)
	if err != nil && !errors.Is(err, pkg.ErrPacienteNoEncontrado) {
		return nil, err
	}

	if existe != nil {
		return nil, pkg.ErrDNIDuplicado
	}

	paciente := &models.Paciente{
		Nombre:          req.Nombre,
		Apellido:        req.Apellido,
		DNI:             req.DNI,
		Telefono:        req.Telefono,
		Email:           req.Email,
		FechaNacimiento: req.FechaNacimiento,
		Direccion:       req.Direccion,
		ObraSocial:      req.ObraSocial,
		Activo:          true,
	}

	if err := s.repoPaciente.CrearPaciente(ctx, paciente); err != nil {
		return nil, err
	}

	return paciente, nil
}

func (s *pacienteService) ObtenerPacientePorID(ctx context.Context, authUserID int64, authRol userModel.Rol, pacienteID int64) (*models.Paciente, error) {

	if pacienteID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(
		ctx,
		pacienteID,
	)
	if err != nil {
		return nil, err
	}

	if !paciente.Activo {
		return nil, pkg.ErrPacienteInactivo
	}

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		// paciente sin medico asignado
		if paciente.MedicoTratante == nil {
			return nil, pkg.ErrForbidden
		}

		medico, err := s.repoMedico.ObtenerMedicoPorID(
			ctx,
			*paciente.MedicoTratante,
		)
		if err != nil {
			return nil, err
		}

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
	}

	return paciente, nil
}

func (s *pacienteService) ObtenerPacientePorDNI(ctx context.Context, dni string) (*models.Paciente, error) {

	if strings.TrimSpace(dni) == "" {
		return nil, pkg.ErrDNIInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorDNI(ctx, dni)
	if err != nil {
		return nil, err
	}

	if !paciente.Activo {
		return nil, pkg.ErrPacienteInactivo
	}
	return paciente, nil
}

func (s *pacienteService) ListarPacientesActivos(ctx context.Context, authUserID int64, authRol userModel.Rol) ([]*models.Paciente, error) {

	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorUserID(
			ctx,
			authUserID,
		)
		if err != nil {
			return nil, err
		}

		pacientes, err := s.repoPaciente.ListarPacientesPorMedico(
			ctx,
			medico.ID,
		)
		if err != nil {
			return nil, err
		}

		if pacientes == nil {
			return []*models.Paciente{}, nil
		}

		return pacientes, nil
	}

	// ADMIN / ADMINISTRATIVO
	pacientes, err := s.repoPaciente.ListarPacientesActivos(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"error al listar pacientes activos: %w",
			err,
		)
	}

	if pacientes == nil {
		return []*models.Paciente{}, nil
	}

	return pacientes, nil
}

func (s *pacienteService) DesactivarPaciente(ctx context.Context, pacienteID int64) error {
	if pacienteID <= 0 {
		return pkg.ErrIDInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return err
	}

	if !paciente.Activo {
		return pkg.ErrPacienteInactivo
	}

	if err := s.repoPaciente.DesactivarPaciente(ctx, pacienteID); err != nil {
		return fmt.Errorf("error al desactivar paciente %w", err)
	}

	return nil
}

func (s *pacienteService) ActivarPaciente(ctx context.Context, pacienteID int64) error {
	if pacienteID <= 0 {
		return pkg.ErrIDInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return err
	}
	if paciente.Activo {
		return pkg.ErrPacienteYaActivo
	}

	if err := s.repoPaciente.ActivarPaciente(ctx, pacienteID); err != nil {
		return fmt.Errorf("activar paciente: %w", err)
	}

	return nil

}

func (s *pacienteService) ActualizarPaciente(ctx context.Context, authUserID int64, authRol userModel.Rol, pacienteID int64, req dto.ActualizarPacienteRequest) (*models.Paciente, error) {
	if pacienteID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return nil, err
	}
	if !paciente.Activo {
		return nil, pkg.ErrPacienteInactivo
	}

	if authRol == userModel.RolMedico {

		if paciente.MedicoTratante == nil {
			return nil, pkg.ErrForbidden
		}

		medico, err := s.repoMedico.ObtenerMedicoPorID(
			ctx,
			*paciente.MedicoTratante,
		)
		if err != nil {
			return nil, err
		}

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
	}

	// Helper local para limpiar la vista
	actualizarSiValido := func(destino *string, valor string) {
		if val := strings.TrimSpace(valor); val != "" {
			*destino = val
		}
	}

	actualizarSiValido(&paciente.Nombre, req.Nombre)
	actualizarSiValido(&paciente.Apellido, req.Apellido)
	actualizarSiValido(&paciente.Telefono, req.Telefono)
	actualizarSiValido(&paciente.Email, req.Email)
	actualizarSiValido(&paciente.Direccion, req.Direccion)
	actualizarSiValido(&paciente.ObraSocial, req.ObraSocial)

	if err := s.repoPaciente.ActualizarPaciente(ctx, paciente); err != nil {
		return nil, fmt.Errorf("error al actualizar paciente %w", err)
	}

	return paciente, nil
}

func (s *pacienteService) AsignarMedicoTratante(ctx context.Context, pacienteID, medicoID int64) error {
	if pacienteID <= 0 || medicoID <= 0 {
		return pkg.ErrIDInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return err
	}
	if !paciente.Activo {
		return pkg.ErrPacienteInactivo
	}

	medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return err
	}
	if !medico.Activo {
		return pkg.ErrMedicoInactivo
	}

	if paciente.MedicoTratante != nil && *paciente.MedicoTratante == medicoID {
		return nil
	}

	if err := s.repoPaciente.AsignarMedicoTratante(ctx, pacienteID, medicoID); err != nil {
		return fmt.Errorf("error al asignar medico: %w", err)
	}

	return nil
}

func (s *pacienteService) QuitarMedicoTratante(ctx context.Context, pacienteID int64) error {
	if pacienteID <= 0 {
		return pkg.ErrIDInvalido
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return err
	}

	if paciente.MedicoTratante == nil {
		return nil
	}

	if !paciente.Activo {
		return pkg.ErrPacienteInactivo
	}

	if err := s.repoPaciente.QuitarMedicoTratante(ctx, pacienteID); err != nil {
		return fmt.Errorf("error al quitar medico del paciente: %w", err)
	}

	return nil
}

func (s *pacienteService) ListarPacientesPorMedico(ctx context.Context, authUserID int64, authRol userModel.Rol, medicoID int64) ([]*models.Paciente, error) {
	if medicoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if !medico.Activo {
		return nil, pkg.ErrMedicoInactivo
	}

	if authRol == userModel.RolMedico {

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
	}

	pacientes, err := s.repoPaciente.ListarPacientesPorMedico(ctx, medicoID)
	if err != nil {
		return nil, fmt.Errorf("error al listar paciente por medico: %w", err)
	}

	return pacientes, nil
}
