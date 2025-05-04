package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/metadata"
)

const (
	// UserIDHeader header name for passing user ID between gRPC services.
	UserIDHeader string = "x-user-id"
	// UserRolesHeader header name for passing user role ids between gRPC services.
	UserRolesHeader string = "x-user-roles"
)

var (
	// ErrNoAuthMetadata is an error when no required auth metadata was found.
	ErrNoAuthMetadata = errors.New("no auth metadata in context")
)

// GetUserMeta retrieves auth metadata from incoming gRPC context.
func GetUserMeta(ctx context.Context) (*authpb.UserAuthMetadata, error) {
	metaRaw, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.WrapErr(ErrNoAuthMetadata)
	}

	userIDRaw, ok := metaRaw[UserIDHeader]
	if !ok {
		return nil, errs.WrapErr(ErrNoAuthMetadata)
	}
	rolesRaw, ok := metaRaw[UserRolesHeader]
	if !ok {
		return nil, errs.WrapErr(ErrNoAuthMetadata)
	}

	userID, err := strconv.ParseInt(userIDRaw[0], 10, 64)
	if err != nil {
		return nil, errs.WrapErr(err, "parse user id")
	}

	rolesSplit := strings.Split(rolesRaw[0], ",")
	roles := make([]int64, len(rolesSplit))
	for idx := range rolesSplit {
		roles[idx], err = strconv.ParseInt(rolesSplit[idx], 10, 64)
		if err != nil {
			return nil, errs.WrapErr(err, "parse role id")
		}
	}

	meta := &authpb.UserAuthMetadata{
		UserId: userID,
		Roles:  roles,
	}

	return meta, nil
}
