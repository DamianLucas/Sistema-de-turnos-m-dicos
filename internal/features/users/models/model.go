package models

import "time"

type Rol string

const (
	RolAdmin          Rol = "admin"
	RolMedico         Rol = "medico"
	RolAdministrativo Rol = "administrativo"
)

type User struct {
	ID       int64  `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Rol      Rol    `json:"rol"`
	Activo   bool   `json:"activo"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
