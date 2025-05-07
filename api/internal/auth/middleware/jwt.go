package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/auth"
	"github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// Jwt is an authorization middleware.
func Jwt(authService pb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken := c.Get("Authorization", "")
		if bearerToken == "" {
			return errs.WrapErr(shared.ErrUnauthorized, "no token in header")
		}

		token := strings.Split(bearerToken, " ")
		if len(token) < 2 {
			return errs.WrapErr(shared.ErrUnauthorized, "invalid token")
		}

		req := &pb.ValidateRequest{Token: token[1]}
		resp, err := authService.Validate(c.UserContext(), req)
		if err != nil {
			return errs.WrapErr(shared.ErrUnauthorized, err.Error())
		}

		meta := resp.GetMeta()
		userID := meta.GetUserId()
		roles := meta.GetRoles()

		c.Locals(shared.UserIDKey, userID)
		c.Locals(shared.UserRolesKey, roles)

		ctx := auth.PushUserMeta(c.UserContext(), &pb.UserAuthMetadata{
			UserId: userID,
			Roles:  roles,
		})

		c.SetUserContext(ctx)

		return c.Next()
	}
}
