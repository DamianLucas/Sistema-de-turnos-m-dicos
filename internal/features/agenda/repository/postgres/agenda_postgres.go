package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"turnos-medicos/internal/features/agenda/models"
	"turnos-medicos/internal/pkg"

	"github.com/lib/pq"
)

type AgendaPostgresRepository struct {
	db *sql.DB
}

func NewAgendaPostgresRepository(db *sql.DB) *AgendaPostgresRepository {
	return &AgendaPostgresRepository{db: db}
}

//Crear metodos de Agenda con sus Query SQL

func (r *AgendaPostgresRepository) CrearAgenda(ctx context.Context, agenda *models.Agenda) error {
	query := `
		INSERT INTO agendas (medico_id, dia_semana, hora_inicio, hora_fin, duracion_turno, activo)
		VALUES ($1, $2, $3, $4, $5, true) 
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		agenda.MedicoID,
		agenda.DiaSemana,
		agenda.HoraInicio,
		agenda.HoraFin,
		agenda.DuracionTurno,
	).Scan(&agenda.ID, &agenda.CreatedAt, &agenda.UpdatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Constraint, "uq_agendas_medico_dia") {
				return pkg.ErrAgendaDuplicada
			}
		}
		return err
	}

	agenda.Activo = true

	return nil

}

func (r *AgendaPostgresRepository) ObtenerAgendaPorID(ctx context.Context, agendaID int64) (*models.Agenda, error) {
	query := `
		SELECT id, medico_id, dia_semana, hora_inicio, hora_fin, duracion_turno, activo, created_at, updated_at 
		FROM agendas
        WHERE id = $1;
	`

	var agenda models.Agenda

	err := r.db.QueryRowContext(ctx, query, agendaID).Scan(
		&agenda.ID,
		&agenda.MedicoID,
		&agenda.DiaSemana,
		&agenda.HoraInicio,
		&agenda.HoraFin,
		&agenda.DuracionTurno,
		&agenda.Activo,
		&agenda.CreatedAt,
		&agenda.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, pkg.ErrAgendaNoEncontrada
	}

	if err != nil {
		return nil, err
	}

	return &agenda, nil

}

func (r *AgendaPostgresRepository) ListarAgendasPorMedico(ctx context.Context, medicoID int64) ([]*models.Agenda, error) {
	query := `
		SELECT id, medico_id, dia_semana, hora_inicio, hora_fin, duracion_turno, activo, created_at, updated_at
		FROM agendas
		WHERE medico_id = $1
		ORDER BY dia_semana;
	`

	rows, err := r.db.QueryContext(ctx, query, medicoID)
	if err != nil {
		return nil, fmt.Errorf("listar agendas por medico: %w", err)
	}
	defer rows.Close()

	agendas := make([]*models.Agenda, 0, 20)

	for rows.Next() {
		agenda := &models.Agenda{}

		err := rows.Scan(
			&agenda.ID,
			&agenda.MedicoID,
			&agenda.DiaSemana,
			&agenda.HoraInicio,
			&agenda.HoraFin,
			&agenda.DuracionTurno,
			&agenda.Activo,
			&agenda.CreatedAt,
			&agenda.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan agenda: %w", err)
		}

		agendas = append(agendas, agenda)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterar agendas: %w", err)
	}

	return agendas, nil
}

func (r *AgendaPostgresRepository) ActualizarAgenda(ctx context.Context, agenda *models.Agenda) error {
	query := `
		UPDATE agendas 
		SET 
			hora_inicio = $1, 
			hora_fin = $2, 
			duracion_turno = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		agenda.HoraInicio,
		agenda.HoraFin,
		agenda.DuracionTurno,
		agenda.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrAgendaNoEncontrada
	}

	return nil
}

func (r *AgendaPostgresRepository) DesactivarAgenda(ctx context.Context, id int64) error {
	query := `
		UPDATE agendas
		SET activo = false,
		    updated_at = NOW()
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrAgendaNoEncontrada
	}

	return nil
}

func (r *AgendaPostgresRepository) ActivarAgenda(ctx context.Context, id int64) error {
	query := `
		UPDATE agendas
		SET activo = true,
		    updated_at = NOW()
		WHERE id = $1;
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ErrAgendaNoEncontrada
	}

	return nil

}
