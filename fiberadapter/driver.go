package fiberadapter

import (
	"fmt"

	"github.com/aziis98/go-authsession"
	"github.com/gofiber/fiber/v2"
)

type fiberDriver struct{}

func (_ *fiberDriver) Cookie(ctx any, name string) (string, error) {
	c, ok := ctx.(*fiber.Ctx)
	if !ok {
		return "", fmt.Errorf(`invalid context, must be of type "*fiber.Ctx"`)
	}

	return c.Cookies(name), nil
}

func (_ *fiberDriver) SetCookie(ctx any, cookie authsession.Cookie) error {
	c, ok := ctx.(*fiber.Ctx)
	if !ok {
		return fmt.Errorf(`invalid context, must be of type "*fiber.Ctx"`)
	}

	c.Cookie(&fiber.Cookie{
		Name:    cookie.Name,
		Path:    cookie.Path,
		Value:   cookie.Value,
		Expires: cookie.Expires,
	})

	return nil
}

func New() authsession.CookieAdapter {
	return &fiberDriver{}
}

func LoggedMiddleware[UserId comparable](service authsession.Service[UserId]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logged, err := service.IsLogged(c)
		if err != nil {
			return err
		}

		if !logged {
			return authsession.ErrNotAuthorized
		}

		return c.Next()
	}
}

func PermissionsMiddleware[UserId comparable](service authsession.Service[UserId], required []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := service.RequestUser(c)
		if err != nil {
			return err
		}

		ok, err := service.Authenticator.HasPermissions(userId, required)
		if err != nil {
			return err
		}
		if !ok {
			return authsession.ErrNotAuthorized
		}

		return c.Next()
	}
}
