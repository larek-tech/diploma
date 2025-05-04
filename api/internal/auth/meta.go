package auth

import (
	"context"
	"strconv"
	"strings"

	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	grpcclient "github.com/yogenyslav/pkg/grpc_client"
)

// PushUserMeta propagates auth metadata into outgoing gRPC context.
func PushUserMeta(ctx context.Context, meta *authpb.UserAuthMetadata) context.Context {
	ctx = grpcclient.PushOutMeta(ctx, shared.UserIDHeader, strconv.FormatInt(meta.GetUserId(), 10))
	rolesRaw := meta.GetRoles()
	roles := make([]string, len(rolesRaw))
	for idx := range roles {
		roles[idx] = strconv.FormatInt(rolesRaw[idx], 10)
	}
	ctx = grpcclient.PushOutMeta(ctx, shared.UserRolesHeader, strings.Join(roles, ","))
	return ctx
}
