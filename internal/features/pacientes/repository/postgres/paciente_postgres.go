package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"turnos-medicos/internal/features/pacientes/models"
	"turnos-medicos/internal/pkg"

	"github.com/lib/pq"
)

type PacientePostgresRepository struct {
	db *sql.DB
}

func NewPacientePostgresRepository(db *sql.DB) *PacientePostgresRepository {
	return &PacientePostgresRepository{db: db}
}

//Crear metodos de Pacientes con sus Query SQL

func (r *PacientePostgresRepository) CrearPaciente(ctx context.Context, p *models.Paciente) error {
	query := `
		INSERT INTO pacientes (
			nombre,
			apellido,
			dni,
			email,
			telefono,
			fecha_nacimiento,
			direccion,
			obra_social,
			medico_tratante_id,
			activo
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, true
		)
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		p.Nombre,
		p.Apellido,
		p.DNI,
		p.Email,
		p.Telefono,
		p.FechaNacimiento,
		p.Direccion,
		p.ObraSocial,
		p.MedicoTratante,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "pacientes_dni_key" {
				return pkg.ErrDNIDuplicado
			}
		}
		return err
	}

	return nil
}

func (r *PacientePostgresRepository) ObtenerPacientePorID(ctx context.Context, id int64) (*models.Paciente, error) {

	query := `
        SELECT
            id,
            nombre,
            apellido,
            dni,
            email,
            telefono,
            fecha_nacimiento,
            direccion,
            obra_social,
            medico_tratante_id,
            activo,
            created_at,
            updated_at
        FROM pacientes
        WHERE id = $1;
    `

	var paciente models.Paciente
	var medicoID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&paciente.ID,
		&paciente.Nombre,
		&paciente.Apellido,
		&paciente.DNI,
		&paciente.Email,
		&paciente.Telefono,
		&paciente.FechaNacimiento,
		&paciente.Direccion,
		&paciente.ObraSocial,
		&medicoID,
		&paciente.Activo,
		&paciente.CreatedAt,
		&paciente.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrPacienteNoEncontrado
		}
		return nil, err
	}

	// manejo de NULL
	if medicoID.Valid {
		paciente.MedicoTratante = &medicoID.Int64
	} else {
		paciente.MedicoTratante = nil
	}

	return &paciente, nil
}

func (r *PacientePostgresRepository) ObtenerPacientePorDNI(ctx context.Context, dni string) (*models.Paciente, error) {

	query := `SELECT id, nombre, apellido, dni, email, telefono, fecha_nacimiento, direccion, obra_social, medico_tratante_id, activo, created_at, updated_at 
		FROM pacientes
		WHERE dni = $1`

	var paciente models.Paciente
	var medicoID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, dni).Scan(
		&paciente.ID,
		&paciente.Nombre,
		&paciente.Apellido,
		&paciente.DNI,
		&paciente.Email,
		&paciente.Telefono,
		&paciente.FechaNacimiento,
		&paciente.Direccion,
		&paciente.ObraSocial,
		&medicoID, //cambio acá
		&paciente.Activo,
		&paciente.CreatedAt,
		&paciente.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrPacienteNoEncontrado
		}
		return nil, err
	}

	// manejo de NULL
	if medicoID.Valid {
		paciente.MedicoTratante = &medicoID.Int64
	} else {
		paciente.MedicoTratante = nil
	}

	return &paciente, nil
}

func (r *PacientePostgresRepository) ListarPacientesActivos(ctx context.Context) ([]*models.Paciente, error) {

	query := `SELECT id, nombre, apellido, dni, email, telefono, fecha_nacimiento, direccion, obra_social, medico_tratante_id, activo, created_at, updated_at 
		FROM pacientes
		WHERE activo = true`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query listar pacientes activos: %w", err)
	}
	defer rows.Close()

	pacientes := make([]*models.Paciente, 0, 30)

	for rows.Next() {
		paciente := &models.Paciente{}
		var medicoID sql.NullInt64

		err := rows.Scan(
			&paciente.ID,
			&paciente.Nombre,
			&paciente.Apellido,
			&paciente.DNI,
			&paciente.Email,
			&paciente.Telefono,
			&paciente.FechaNacimiento,
			&paciente.Direccion,
			&paciente.ObraSocial,
			&medicoID, // cambio acá
			&paciente.Activo,
			&paciente.CreatedAt,
			&paciente.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando paciente: %w", err)
		}

		//manejo de NULL
		if medicoID.Valid {
			paciente.MedicoTratante = &medicoID.Int64
		} else {
			paciente.MedicoTratante = nil
		}

		pacientes = append(pacientes, paciente)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando pacientes: %w", err)
	}

	return pacientes, nil
}

func (r *PacientePostgresRepository) ActualizarPaciente(ctx context.Context, p *models.Paciente) error {
	query := `
		UPDATE pacientes
		SET
			nombre = $1,
			apellido = $2,
			dni = $3,
			email = $4,
			telefono = $5,
			fecha_nacimiento = $6,
			direccion = $7,
			obra_social = $8,
			medico_tratante_id = $9,
			updated_at = NOW()
		WHERE id = $10
		RETURNING
			id,
			nombre,
			apellido,
			dni,
			email,
			telefono,
			fecha_nacimiento,
			direccion,
			obra_social,
			medico_tratante_id,
			activo,
			created_at,
			updated_at;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		p.Nombre,
		p.Apellido,
		p.DNI,
		p.Email,
		p.Telefono,
		p.FechaNacimiento,
		p.Direccion,
		p.ObraSocial,
		p.MedicoTratante,
		p.ID,
	).Scan(
		&p.ID,
		&p.Nombre,
		&p.Apellido,
		&p.DNI,
		&p.Email,
		&p.Telefono,
		&p.FechaNacimiento,
		&p.Direccion,
		&p.ObraSocial,
		&p.MedicoTratante,
		&p.Activo,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {

		// 🔥 paciente no encontrado
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrPacienteNoEncontrado
		}

		// 🔥 DNI duplicado
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "pacientes_dni_key" {
				return pkg.ErrDNIDuplicado
			}
		}

		return err
	}

	return nil
}

func (r *PacientePostgresRepository) DesactivarPaciente(ctx context.Context, pacienteID int64) error {

	query := `
		UPDATE pacientes
		SET activo = false,
		    updated_at = NOW()
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, pacienteID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrPacienteNoEncontrado
	}

	return nil
}

func (r *PacientePostgresRepository) AsignarMedicoTratante(ctx context.Context, pacienteID, medicoID int64) error {
	query := `
		UPDATE pacientes
		SET medico_tratante_id = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`

	resultado, err := r.db.ExecContext(ctx, query, medicoID, pacienteID)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" { // foreign_key_violation
				return pkg.ErrMedicoNoEncontrado
			}
		}
		return err
	}
	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrPacienteNoEncontrado
	}

	return nil
}

func (r *PacientePostgresRepository) QuitarMedicoTratante(ctx context.Context, pacienteID int64) error {
	query := `
		UPDATE pacientes
		SET medico_tratante_id = NULL,
		    updated_at = NOW()
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, pacienteID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrPacienteNoEncontrado
	}

	return nil
}

func (r *PacientePostgresRepository) ListarPacientesPorMedico(ctx context.Context, medicoID int64) ([]*models.Paciente, error) {
	query := `SELECT id, nombre, apellido, dni, email, telefono, fecha_nacimiento, direccion, obra_social, medico_tratante_id, activo, created_at, updated_at
		FROM pacientes
		WHERE medico_tratante_id = $1 AND activo = true`

	rows, err := r.db.QueryContext(ctx, query, medicoID)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query listar pacientes por medico: %w", err)
	}
	defer rows.Close()

	pacientes := make([]*models.Paciente, 0, 20)

	for rows.Next() {
		paciente := &models.Paciente{}
		var medicoIDDB sql.NullInt64

		err := rows.Scan(
			&paciente.ID,
			&paciente.Nombre,
			&paciente.Apellido,
			&paciente.DNI,
			&paciente.Email,
			&paciente.Telefono,
			&paciente.FechaNacimiento,
			&paciente.Direccion,
			&paciente.ObraSocial,
			&medicoIDDB,
			&paciente.Activo,
			&paciente.CreatedAt,
			&paciente.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando paciente: %w", err)
		}

		// aunque no debería ser NULL, mantenemos consistencia
		if medicoIDDB.Valid {
			paciente.MedicoTratante = &medicoIDDB.Int64
		} else {
			paciente.MedicoTratante = nil
		}

		pacientes = append(pacientes, paciente)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando pacientes: %w", err)
	}

	return pacientes, nil
}
