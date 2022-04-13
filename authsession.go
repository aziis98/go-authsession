package authsession

import (
	"errors"
	"log"
	"time"

	"github.com/aziis98/go-authsession/generics"
)

var (
	// ErrInvalidSession is raised when the given session is not present in the session store
	ErrInvalidSession = errors.New(`invalid session`)

	// ErrNotAuthorized is raised when a user tries to access a resource it is not supposed to
	ErrNotAuthorized = errors.New(`not authorized`)
)

// Service is the main service provided by this library. This is generic over "UserId" that is the type of the user ids and "Ctx" that is the context provided by
type Service[UserId comparable] struct {
	Config        ServiceConfig
	Authenticator Authenticator[UserId]
	SessionStore  SessionStore[UserId]
	CookieAdapter CookieAdapter
}

func (service Service[UserId]) Login(ctx any, username, password string) error {
	ok, err := service.Authenticator.Authenticate(username, password)
	if err != nil {
		log.Printf(`err from service.authenticator.Authenticate: %v`, err)
		return err
	}
	if !ok {
		return ErrNotAuthorized
	}

	userId, err := service.Authenticator.FromUsername(username)
	if err != nil {
		return err
	}

	sid, err := service.SessionStore.CreateSession(userId)
	if err != nil {
		log.Printf(`err from service.sessionStore.newSession: %v`, err)
		return err
	}

	service.CookieAdapter.SetCookie(ctx, Cookie{
		Name:    service.Config.CookieName,
		Path:    service.Config.CookiePath,
		Value:   sid,
		Expires: time.Now().Add(service.Config.CookieDuration),
	})

	return nil
}

func (service Service[UserId]) IsLogged(ctx any) (bool, error) {
	sid, err := service.CookieAdapter.Cookie(ctx, service.Config.CookieName)
	if err != nil {
		return false, err
	}
	if sid == "" {
		return false, nil
	}

	if _, err := service.SessionStore.UserForSession(sid); err != nil {
		return false, err
	}

	return true, nil
}

func (service Service[UserId]) RequestUser(ctx any) (UserId, error) {
	sid, err := service.CookieAdapter.Cookie(ctx, service.Config.CookieName)
	if err != nil {
		return generics.Zero[UserId](), err
	}
	if sid == "" {
		return generics.Zero[UserId](), nil
	}

	userId, err := service.SessionStore.UserForSession(sid)
	if err != nil {
		return generics.Zero[UserId](), err
	}

	return userId, nil
}

func (service Service[UserId]) Logout(ctx any) error {
	sid, err := service.CookieAdapter.Cookie(ctx, service.Config.CookieName)
	if err != nil {
		return err
	}

	if err := service.SessionStore.DeleteSession(sid); err != nil {
		return err
	}

	service.CookieAdapter.SetCookie(ctx, Cookie{
		Name:    service.Config.CookieName,
		Path:    service.Config.CookiePath,
		Value:   "",
		Expires: time.Now(),
	})

	return nil
}
