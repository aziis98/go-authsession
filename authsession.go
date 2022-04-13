package authsession

import (
	"errors"

	"github.com/aziis98/go-authsession/session"
)

var (
	// ErrNotAuthorized is raised when a user tries to access a resource it is not supposed to
	ErrNotAuthorized = errors.New(`not authorized`)

	// ErrNoPermissionChecker is raised when one tries to use "(*Basic).HasPermissions()" on an instance without a valid PermissionChecker
	ErrNoPermissionChecker = errors.New(`no permission checker`)
)

// Base is the main Base provided by this library
type Base struct {
	Config *CookieConfig

	credChecker  CredentialChecker
	permChecker  PermissionChecker
	sessionStore session.Store
}

func (s *Base) Login(userId string, password string) (string, error) {
	ok, err := s.credChecker.CheckCredentials(userId, password)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrNotAuthorized
	}

	sid, err := s.sessionStore.CreateSession(userId)
	if err != nil {
		return "", err
	}

	return sid, nil
}

func (s *Base) HasPermissions(userId string, required []string) (bool, error) {
	if s.permChecker == nil {
		panic(ErrNoPermissionChecker)
	}

	ok, err := s.permChecker.HasPermissions(userId, required)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *Base) IsLogged(sessionId string) (bool, error) {
	// shortcut for passing directly the cookie value from http framework of choice
	if sessionId == "" {
		return false, nil
	}

	if _, err := s.sessionStore.UserForSession(sessionId); err != nil {
		return false, err
	}

	return true, nil
}

func (service *Base) UserForSession(sessionId string) (string, error) {
	// shortcut for passing directly the cookie value from http framework of choice
	if sessionId == "" {
		return "", session.ErrSessionNotFound
	}

	userId, err := service.sessionStore.UserForSession(sessionId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (service *Base) Logout(sessionId string) error {
	// shortcut for passing directly the cookie value from http framework of choice
	if sessionId == "" {
		return session.ErrSessionNotFound
	}

	if err := service.sessionStore.DeleteSession(sessionId); err != nil {
		return err
	}

	return nil
}
