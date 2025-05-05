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
