package shared

type contextKey uint8

const (
	// UserIDKey is a key for storing userID in context.
	UserIDKey contextKey = iota
	// UserRolesKey is a key for storing user roles in context.
	UserRolesKey
)
