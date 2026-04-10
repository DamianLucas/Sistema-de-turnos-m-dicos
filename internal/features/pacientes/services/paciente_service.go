package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"turnos-medicos/internal/features/pacientes/dto"
	"turnos-medicos/internal/features/pacientes/models"
	"turnos-medicos/internal/pkg"

	repositoryMedico "turnos-medicos/internal/features/medicos/repository"
	repositoryPaciente "turnos-medicos/internal/features/pacientes/repository"
)

type PacienteService interface {
	CrearPaciente(ctx context.Context, req dto.CrearPacienteRequest) (*models.Paciente, error)
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

func (s *pacienteService) ObtenerPacientePorID(ctx context.Context, pacienteID int64) (*models.Paciente, error) {
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

func (s *pacienteService) ListarPacientesActivos(ctx context.Context) ([]*models.Paciente, error) {
	pacientes, err := s.repoPaciente.ListarPacientesActivos(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", pkg.ErrListarPacientesActivos, err)
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
		return fmt.Errorf("%w: %v", pkg.ErrDesactivarPaciente, err)
	}

	return nil
}

func (s *pacienteService) ActualizarPaciente(ctx context.Context, pacienteID int64, req dto.ActualizarPacienteRequest) (*models.Paciente, error) {
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
		return nil, fmt.Errorf("%w: %v", pkg.ErrActualizarPaciente, err)
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

	paciente.MedicoTratante = &medicoID

	if err := s.repoPaciente.AsignarMedicoTratante(ctx, pacienteID, medicoID); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrAsignarMedicoPaciente, err)
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

	if !paciente.Activo {
		return pkg.ErrPacienteInactivo
	}

	if paciente.MedicoTratante == nil {
		return nil
	}

	if err := s.repoPaciente.QuitarMedicoTratante(ctx, pacienteID); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrQuitarMedicoPaciente, err)
	}

	return nil
}

func (s *pacienteService) ListarPacientesPorMedico(ctx context.Context, medicoID int64) ([]*models.Paciente, error) {
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

	pacientes, err := s.repoPaciente.ListarPacientesPorMedico(ctx, medicoID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", pkg.ErrListarPacientesPorMedico, err)
	}

	return pacientes, nil
}
