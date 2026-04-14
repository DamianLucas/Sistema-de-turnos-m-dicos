package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"turnos-medicos/internal/features/medicos/dto"
	"turnos-medicos/internal/pkg"

	medicoModel "turnos-medicos/internal/features/medicos/models"
	userModel "turnos-medicos/internal/features/users/models"

	medicoRepo "turnos-medicos/internal/features/medicos/repository"
	pacienteRepo "turnos-medicos/internal/features/pacientes/repository"
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
	ActivarMedico(ctx context.Context, medicoID int64) error
}

type medicoService struct {
	medicoRepo   medicoRepo.MedicoRepository
	userRepo     userRepo.UserRepository
	pacienteRepo pacienteRepo.PacienteRepository
	db           *sql.DB
}

func NewMedicoService(medicoRepo medicoRepo.MedicoRepository, userRepo userRepo.UserRepository, pacienteRepo pacienteRepo.PacienteRepository, db *sql.DB) MedicoService {
	return &medicoService{
		medicoRepo:   medicoRepo,
		userRepo:     userRepo,
		pacienteRepo: pacienteRepo,
		db:           db,
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
	medico.UserID = user.ID

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
		return nil, fmt.Errorf("%w: %w", pkg.ErrListarMedicosActivos, err)
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

	if req.Especialidad != "" {
		medicoActual.Especialidad = req.Especialidad
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

	if medicoID <= 1 {
		return pkg.ErrIDInvalido
	}

	medico, err := s.medicoRepo.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return err
	}

	if !medico.Activo {
		return pkg.ErrMedicoInactivo
	}

	// 🔥 BEGIN TX
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}

	// 🔥 manejo seguro de rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 🔥 1. remover médico de pacientes
	queryPacientes := `
		UPDATE pacientes 
		SET medico_tratante_id = NULL,
		    updated_at = NOW()
		WHERE medico_tratante_id = $1;
	`

	if _, err = tx.ExecContext(ctx, queryPacientes, medicoID); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrQuitarMedicoPaciente, err)
	}

	// 🔥 2. desactivar médico
	queryUser := `
	UPDATE users
	SET activo = false,
	    updated_at = NOW()
	WHERE id = $1;
`

	result, err := tx.ExecContext(ctx, queryUser, medico.UserID)
	if err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrDesactivarUsuario, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return pkg.ErrUsuarioNoEncontrado
	}

	// 🔥 COMMIT
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al hacer commit: %w", err)
	}

	return nil
}

func (s *medicoService) ActivarMedico(ctx context.Context, medicoID int64) error {
	return s.medicoRepo.ActivarMedico(ctx, medicoID)
}
