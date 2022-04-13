package authsession

import (
	"errors"

	"github.com/aziis98/go-authsession/session"
)

var (
	// ErrUserNotFound should be raised when a credential checker is called with a non-existant user
	ErrUserNotFound = errors.New(`user not found`)

	// ErrUnauthorized is raised when a user tries to access a resource it is not supposed to
	ErrUnauthorized = errors.New(`not authorized`)

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

// Login the given user if the provided credentials are correct this returns a new session token associated with this user.
func (b *Base) Login(userId string, password string) (string, error) {
	ok, err := b.credChecker.CheckCredentials(userId, password)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrUnauthorized
	}

	sid, err := b.sessionStore.CreateSession(userId)
	if err != nil {
		return "", err
	}

	return sid, nil
}

// Logout deletes the given session if still valid.
func (b *Base) Logout(sessionId string) error {
	// shortcut for passing directly the cookie value from http framework of choice
	if sessionId == "" {
		return session.ErrSessionNotFound
	}

	if err := b.sessionStore.DeleteSession(sessionId); err != nil {
		return err
	}

	return nil
}

// UserForSession retrieves the user associated with this session.
func (b *Base) UserForSession(sessionId string) (string, error) {
	// shortcut for passing directly the cookie value from http framework of choice
	if sessionId == "" {
		return "", session.ErrSessionNotFound
	}

	userId, err := b.sessionStore.UserForSession(sessionId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

// IsLogged returns true if the given session is valid and associated with a user
func (b *Base) IsLogged(sessionId string) (bool, error) {
	if _, err := b.sessionStore.UserForSession(sessionId); err != nil {
		return false, err
	}

	return true, nil
}

func (b *Base) HasPermissions(userId string, required []string) (bool, error) {
	if b.permChecker == nil {
		panic(ErrNoPermissionChecker)
	}

	ok, err := b.permChecker.HasPermissions(userId, required)
	if err != nil {
		return false, err
	}

	return ok, nil
}
