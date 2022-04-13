package fiberauth

import (
	"time"

	"github.com/aziis98/go-authsession"
	"github.com/gofiber/fiber/v2"
)

type FiberAuth struct {
	basic *authsession.Base
}

func (fa *FiberAuth) cookie(c *fiber.Ctx) (string, error) {
	return c.Cookies(fa.basic.Config.Name), nil
}

func (fa *FiberAuth) setCookie(c *fiber.Ctx, value string) {
	c.Cookie(&fiber.Cookie{
		Name:    fa.basic.Config.Name,
		Path:    fa.basic.Config.Path,
		Value:   value,
		Expires: time.Now().Add(fa.basic.Config.Duration),
	})
}

func (fa *FiberAuth) clearCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    fa.basic.Config.Name,
		Path:    fa.basic.Config.Path,
		Value:   "",
		Expires: time.Now(),
	})
}

func (fa *FiberAuth) Login(c *fiber.Ctx, username, password string) error {
	newSessionId, err := fa.basic.Login(username, password)
	if err != nil {
		return err
	}

	fa.setCookie(c, newSessionId)

	return nil
}

func (fa *FiberAuth) Logout(c *fiber.Ctx) error {
	sessionId, err := fa.cookie(c)
	if err != nil {
		return err
	}

	if _, err := fa.basic.UserForSession(sessionId); err != nil {
		return err
	}

	fa.clearCookie(c)

	return nil
}

func (fa *FiberAuth) RequestUser(c *fiber.Ctx) (string, error) {
	sessionId, err := fa.cookie(c)
	if err != nil {
		return "", err
	}

	userId, err := fa.basic.UserForSession(sessionId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (fa *FiberAuth) IsLogged(c *fiber.Ctx) (bool, error) {
	if _, err := fa.RequestUser(c); err != nil {
		return false, err
	}

	return true, nil
}

func (fa *FiberAuth) LoggedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logged, err := fa.IsLogged(c)
		if err != nil {
			return err
		}

		if !logged {
			return authsession.ErrNotAuthorized
		}

		return c.Next()
	}
}

func (fa *FiberAuth) PermissionsMiddleware(required []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := fa.RequestUser(c)
		if err != nil {
			return err
		}

		ok, err := fa.basic.HasPermissions(userId, required)
		if err != nil {
			return err
		}

		if !ok {
			return authsession.ErrNotAuthorized
		}

		return c.Next()
	}
}
