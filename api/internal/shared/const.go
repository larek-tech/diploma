package shared

type contextKey uint8

const (
	// UserIDKey is a key for storing userID in context.
	UserIDKey contextKey = iota
	// UserRolesKey is a key for storing user roles in context.
	UserRolesKey
)

const (
	// UserIDHeader header name for passing user ID between gRPC services.
	UserIDHeader string = "x-user-id"
	// UserRolesHeader header name for passing user role ids between gRPC services.
	UserRolesHeader string = "x-user-roles"
)
