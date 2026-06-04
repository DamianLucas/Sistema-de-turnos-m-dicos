package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"turnos-medicos/internal/features/turnos/models"
	"turnos-medicos/internal/pkg"

	"github.com/lib/pq"
)

type TurnoPostgresRepository struct {
	db *sql.DB
}

func NewTurnoPostgresRepository(db *sql.DB) *TurnoPostgresRepository {
	return &TurnoPostgresRepository{db: db}
}

//Crear metodos de Medicos con sus Query SQL

func (r *TurnoPostgresRepository) CrearTurno(ctx context.Context, turno *models.Turno) error {
	query := `
			INSERT INTO turnos (
			agenda_id,
			medico_id,
			paciente_id,
			fecha,
			hora_inicio,
			hora_fin,
			estado
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		turno.AgendaID,
		turno.MedicoID,
		turno.PacienteID,
		turno.Fecha,
		turno.HoraInicio,
		turno.HoraFin,
		turno.Estado,
	).Scan(
		&turno.ID,
		&turno.CreatedAt,
		&turno.UpdatedAt,
	)

	if err != nil {

		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {

			if strings.Contains(
				pqErr.Constraint,
				"uq_turnos_medico_fecha_hora",
			) {
				return pkg.ErrTurnoDuplicado
			}
		}
		return err
	}

	return nil
}

func (r *TurnoPostgresRepository) ObtenerTurnoPorID(ctx context.Context, turnoID int64) (*models.Turno, error) {
	query := `
		SELECT 
			id,
			agenda_id,
			medico_id,
			paciente_id,
			fecha,
			hora_inicio,
			hora_fin,
			estado,
			created_at,
			updated_at
		FROM turnos
        WHERE id = $1;
	`
	var turno models.Turno

	err := r.db.QueryRowContext(ctx, query, turnoID).Scan(
		&turno.ID,
		&turno.AgendaID,
		&turno.MedicoID,
		&turno.PacienteID,
		&turno.Fecha,
		&turno.HoraInicio,
		&turno.HoraFin,
		&turno.Estado,
		&turno.CreatedAt,
		&turno.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkg.ErrTurnoNoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &turno, nil

}

func (r *TurnoPostgresRepository) ListarTurnosPorMedico(ctx context.Context, medicoID int64) ([]*models.Turno, error) {
	query := `
		SELECT
			id,
			agenda_id,
			medico_id,
			paciente_id,
			fecha,
			hora_inicio,
			hora_fin,
			estado,
			created_at,
			updated_at
		FROM turnos
		WHERE medico_id = $1
		ORDER BY fecha, hora_inicio;
	`
	rows, err := r.db.QueryContext(ctx, query, medicoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	turnos := make([]*models.Turno, 0, 20)

	for rows.Next() {
		var turno models.Turno

		err := rows.Scan(
			&turno.ID,
			&turno.AgendaID,
			&turno.MedicoID,
			&turno.PacienteID,
			&turno.Fecha,
			&turno.HoraInicio,
			&turno.HoraFin,
			&turno.Estado,
			&turno.CreatedAt,
			&turno.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		turnos = append(turnos, &turno)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return turnos, nil
}

func (r *TurnoPostgresRepository) ListarTurnosPorPaciente(ctx context.Context, pacienteID int64) ([]*models.Turno, error) {
	query := `
		SELECT
			id,
			agenda_id,
			medico_id,
			paciente_id,
			fecha,
			hora_inicio,
			hora_fin,
			estado,
			created_at,
			updated_at
		FROM turnos
		WHERE paciente_id = $1
		ORDER BY fecha, hora_inicio;
	`

	rows, err := r.db.QueryContext(ctx, query, pacienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	turnos := make([]*models.Turno, 0, 20)

	for rows.Next() {
		var turno models.Turno

		err := rows.Scan(
			&turno.ID,
			&turno.AgendaID,
			&turno.MedicoID,
			&turno.PacienteID,
			&turno.Fecha,
			&turno.HoraInicio,
			&turno.HoraFin,
			&turno.Estado,
			&turno.CreatedAt,
			&turno.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		turnos = append(turnos, &turno)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return turnos, nil

}

func (r *TurnoPostgresRepository) ListarTurnosDisponibles(ctx context.Context, medicoID int64) ([]*models.Turno, error) {
	query := `
		SELECT
			id,
			agenda_id,
			medico_id,
			paciente_id,
			fecha,
			hora_inicio,
			hora_fin,
			estado,
			created_at,
			updated_at
		FROM turnos
		WHERE medico_id = $1
		  AND estado = 'disponible'
		ORDER BY fecha, hora_inicio;
	`

	rows, err := r.db.QueryContext(ctx, query, medicoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	turnos := make([]*models.Turno, 0, 20)

	for rows.Next() {
		var turno models.Turno

		err := rows.Scan(
			&turno.ID,
			&turno.AgendaID,
			&turno.MedicoID,
			&turno.PacienteID,
			&turno.Fecha,
			&turno.HoraInicio,
			&turno.HoraFin,
			&turno.Estado,
			&turno.CreatedAt,
			&turno.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		turnos = append(turnos, &turno)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return turnos, nil

}

func (r *TurnoPostgresRepository) ReservarTurno(ctx context.Context, turnoID int64, pacienteID int64) error {
	query := `
		UPDATE turnos
		SET
			paciente_id = $1,
			estado = 'reservado',
			updated_at = NOW()
		WHERE id = $2;
	`

	resultado, err := r.db.ExecContext(
		ctx,
		query,
		pacienteID,
		turnoID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrTurnoNoEncontrado
	}

	return nil
}

func (r *TurnoPostgresRepository) LiberarTurno(ctx context.Context, turnoID int64) error {
	query := `
		UPDATE turnos
		SET
			paciente_id = NULL,
			estado = 'disponible',
			updated_at = NOW()
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, turnoID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrTurnoNoEncontrado
	}

	return nil
}

func (r *TurnoPostgresRepository) MarcarTurnoAtendido(ctx context.Context, turnoID int64) error {
	query := `
		UPDATE turnos
		SET
			estado = 'atendido',
			updated_at = NOW()
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, turnoID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrTurnoNoEncontrado
	}

	return nil
}

func (r *TurnoPostgresRepository) MarcarTurnoNoAsistio(ctx context.Context, turnoID int64) error {
	query := `
		UPDATE turnos
		SET
			estado = 'no_asistio',
			updated_at = NOW()
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, turnoID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrTurnoNoEncontrado
	}

	return nil
}
