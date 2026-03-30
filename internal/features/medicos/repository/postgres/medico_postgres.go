package postgres

import (
	"context"
	"database/sql"
	"fmt"
	medicoModel "turnos-medicos/internal/features/medicos/models"
	userModel "turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/pkg"

	"github.com/lib/pq"
)

type MedicoPostgresRepository struct {
	db *sql.DB
}

func NewMedicoPostgresRepository(db *sql.DB) *MedicoPostgresRepository {
	return &MedicoPostgresRepository{db: db}
}

//Crear metodos de Medicos con sus Query SQL

func (r *MedicoPostgresRepository) CrearMedico(ctx context.Context, u *userModel.User, m *medicoModel.Medico) error {

	//Todo pertenece a una misma transaccion
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}() //si el commit no se ejecutó, hace rollback

	queryUser := `INSERT INTO users (nombre, apellido, email, password, rol, activo) 
				  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = tx.QueryRowContext(ctx,
		queryUser,
		u.Nombre,
		u.Apellido,
		u.Email,
		u.Password,
		u.Rol,
		u.Activo,
	).Scan(&u.ID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				switch pqErr.Constraint {
				case "users_email_key":
					return pkg.ErrEmailDuplicado
				default:
					return pkg.ErrUsuarioYaExiste
				}
			}
		}

		return fmt.Errorf("error insertando user: %w", err)
	}

	queryMedico := `INSERT INTO medicos (user_id, matricula, especialidad) VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRowContext(ctx,
		queryMedico,
		u.ID,
		m.Matricula,
		m.Especialidad,
	).Scan(&m.ID)

	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				switch pqErr.Constraint {
				case "medicos_matricula_key":
					return pkg.ErrMatriculaDuplicada
				case "medicos_user_id_key":
					return pkg.ErrUsuarioYaExiste
				}
			}
		}

		return fmt.Errorf("error insertando medico: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error haciendo commit: %w", err)
	}

	return nil
}

func (r *MedicoPostgresRepository) ObtenerMedicoPorID(ctx context.Context, medicoID int64) (*medicoModel.Medico, error) {
	query := `SELECT m.id, m.user_id, u.nombre, u.apellido, u.email, u.activo, m.matricula, m.especialidad, created_at, updated_at
			FROM medicos m
			INNER JOIN users u ON u.id = m.user_id
			WHERE m.id = $1`

	var medico medicoModel.Medico

	err := r.db.QueryRowContext(ctx, query, medicoID).Scan(
		&medico.ID,
		&medico.UserID,
		&medico.Nombre,
		&medico.Apellido,
		&medico.Email,
		&medico.Activo,
		&medico.Matricula,
		&medico.Especialidad,
		&medico.CreatedAt,
		&medico.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkg.ErrMedicoNoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &medico, nil

}

func (r *MedicoPostgresRepository) ObtenerMedicoPorMatricula(ctx context.Context, matricula string) (*medicoModel.Medico, error) {
	query := `SELECT m.id, m.user_id, u.nombre, u.apellido, u.email, u.activo, m.matricula, m.especialidad, created_at, updated_at
			FROM medicos m
			INNER JOIN users u ON u.id = m.user_id
			WHERE m.matricula = $1`

	var medico medicoModel.Medico

	err := r.db.QueryRowContext(ctx, query, matricula).Scan(
		&medico.ID,
		&medico.UserID,
		&medico.Nombre,
		&medico.Apellido,
		&medico.Email,
		&medico.Activo,
		&medico.Matricula,
		&medico.Especialidad,
		&medico.CreatedAt,
		&medico.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, pkg.ErrMedicoNoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &medico, nil
}

func (r *MedicoPostgresRepository) ListarMedicosActivos(ctx context.Context) ([]*medicoModel.Medico, error) {
	query := `
	SELECT
		m.id, m.user_id, m.matricula, m.especialidad, m.created_at, m.updated_at,
		u.nombre, u.apellido, u.email, u.activo
	FROM medicos m
	INNER JOIN users u ON u.id = m.user_id
	WHERE u.activo = true	
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query listar medicos activos: %w", err)
	}
	defer rows.Close()

	medicos := make([]*medicoModel.Medico, 0, 20)

	for rows.Next() {
		medico := &medicoModel.Medico{}

		err := rows.Scan(
			&medico.ID,
			&medico.UserID,
			&medico.Matricula,
			&medico.Especialidad,
			&medico.CreatedAt,
			&medico.UpdatedAt,
			&medico.Nombre,
			&medico.Apellido,
			&medico.Email,
			&medico.Activo,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando medico: %w", err)
		}

		medicos = append(medicos, medico)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando medicos %w", err)
	}
	return medicos, nil
}

func (r *MedicoPostgresRepository) ListarMedicosPorEspecialidad(ctx context.Context, especialidad string) ([]*medicoModel.Medico, error) {
	query := `
	SELECT
		m.id, m.user_id, m.matricula, m.especialidad, m.created_at, m.updated_at,
		u.nombre, u.apellido, u.email, u.activo
	FROM medicos m
	INNER JOIN users u ON u.id = m.user_id
	WHERE m.especialidad = $1
	`

	rows, err := r.db.QueryContext(ctx, query, especialidad)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query listar medicos por especialidad: %w", err)
	}
	defer rows.Close()

	medicos := make([]*medicoModel.Medico, 0, 20)

	for rows.Next() {
		medico := &medicoModel.Medico{}

		err := rows.Scan(
			&medico.ID,
			&medico.UserID,
			&medico.Matricula,
			&medico.Especialidad,
			&medico.CreatedAt,
			&medico.UpdatedAt,
			&medico.Nombre,
			&medico.Apellido,
			&medico.Email,
			&medico.Activo,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando medico: %w", err)
		}

		medicos = append(medicos, medico)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando medicos: %w", err)
	}

	return medicos, nil

}

func (r *MedicoPostgresRepository) ActualizarMedico(ctx context.Context, m *medicoModel.Medico) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	//actualiza users
	queryUser := `
	UPDATE users
	SET nombre = $1,
		apellido = $2,
		email = $3,
		updated_at = NOW()
	WHERE id = $4
	`

	res, err := tx.ExecContext(ctx, queryUser,
		m.Nombre,
		m.Apellido,
		m.Email,
		m.UserID,
	)
	if err != nil {
		return fmt.Errorf("error actualizando user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas user: %w", err)
	}
	if rowsAffected == 0 {
		return pkg.ErrUsuarioNoEncontrado
	}

	//actualizada medico
	queryMedico := `
	UPDATE medicos
	SET especialidad = $1,
	    updated_at = NOW()
	WHERE id = $2
	`

	res, err = tx.ExecContext(ctx, queryMedico,
		m.Especialidad,
		m.ID,
	)
	if err != nil {
		return fmt.Errorf("error actualizando medico: %w", err)
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas medico: %w", err)
	}
	if rowsAffected == 0 {
		return pkg.ErrMedicoNoEncontrado
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error haciendo commit: %w", err)
	}

	return nil

}

func (r *MedicoPostgresRepository) DesactivarMedico(ctx context.Context, medicoID int64) error {
	query := `UPDATE users SET activo = false, updated_at = NOW() 
			WHERE id = (SELECT user_id FROM medicos WHERE id = $1)`

	resultado, err := r.db.ExecContext(ctx, query, medicoID)

	if err != nil {
		return err
	}

	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return pkg.ErrMedicoNoEncontrado
	}

	return nil
}
