package model

import "time"

// UserDao is a data layer model.proto for user.
type UserDao struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	HashPassword string    `db:"hash_password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	IsDeleted    bool      `db:"is_deleted"`
}

// Role is a type for user role enum.
type Role string

const (
	// DefaultRole is a role that is assigned to every user on register.
	DefaultRole Role = "default"
	// AdminRole is a role for admin users.
	AdminRole Role = "admin"
)

// RoleDao is a data model.proto for role.
type RoleDao struct {
	ID        int64     `db:"id"`
	Name      Role      `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}
