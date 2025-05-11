package model

import (
	"time"
)

// UserDao is a data model for user on data layer.
type UserDao struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	HashPassword string    `db:"hash_password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	IsDeleted    bool      `db:"is_deleted"`
}
