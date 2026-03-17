package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/utils"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

//Crear metodos de Users con sus Query SQL

// CrearUsuario
func (r *UserPostgresRepository) CrearUsuario(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (nombre, apellido, email, password, rol, activo)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		u.Nombre,
		u.Apellido,
		u.Email,
		u.Password,
		u.Rol,
		u.Activo,
	).Scan(&u.ID)

	return err

}

// ListarUsuarios
func (r *UserPostgresRepository) ListarUsuariosActivos(ctx context.Context) ([]*models.User, error) {

	query := `SELECT id, nombre, apellido, email, rol, activo, created_at, updated_at
	          FROM users
	          WHERE activo = true
	          ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando usuarios activos: %w", err)
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User

		if err := rows.Scan(
			&user.ID,
			&user.Nombre,
			&user.Apellido,
			&user.Email,
			&user.Rol,
			&user.Activo,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ObtenerUsuarioPorID
func (r *UserPostgresRepository) ObtenerUsuarioPorID(ctx context.Context, userID int64) (*models.User, error) {
	query := `SELECT id, nombre, apellido, email, rol, activo, created_at, updated_at FROM users WHERE id = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Nombre,
		&user.Apellido,
		&user.Email,
		&user.Rol,
		&user.Activo,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.ErrUsuarioNoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &user, err
}

// ObtenerUsuarioPorEmail
func (r *UserPostgresRepository) ObtenerUsuarioPorEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, nombre, apellido, email, rol, activo, created_at, updated_at FROM users WHERE email = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Nombre,
		&user.Apellido,
		&user.Email,
		&user.Rol,
		&user.Activo,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.ErrUsuarioNoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &user, err
}

// ObtenerUsuarioPorRol
func (r *UserPostgresRepository) ObtenerUsuarioPorRol(ctx context.Context, userRol models.Rol) ([]*models.User, error) {

	query := `
        SELECT id, nombre, apellido, email, rol, activo, created_at, updated_at
        FROM users
        WHERE rol = $1
    `

	rows, err := r.db.QueryContext(ctx, query, userRol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Nombre,
			&user.Apellido,
			&user.Email,
			&user.Rol,
			&user.Activo,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ActualizarUsuario
func (r *UserPostgresRepository) ActualizarUsuario(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET nombre = $2, apellido = $3, email = $4, password = $5, rol = $6, activo = $7, updated_at = NOW() WHERE id = $1`

	resultado, err := r.db.ExecContext(
		ctx,
		query,
		u.ID,
		u.Nombre,
		u.Apellido,
		u.Email,
		u.Password,
		u.Rol,
		u.Activo,
	)

	if err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}
	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DesactivarUsuario
func (r *UserPostgresRepository) DesactivarUsuario(ctx context.Context, userID int64) error {
	query := `UPDATE users SET activo = false, updated_at = NOW() WHERE id = $1`

	resultado, err := r.db.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return utils.ErrUsuarioNoEncontrado
	}

	return nil
}
