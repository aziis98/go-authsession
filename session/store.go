package session

import "errors"

var (
	// ErrSessionNotFound is raised when the given session is not present in the session store
	ErrSessionNotFound = errors.New(`session not found`)

	// ErrSessionInvalid is raised when the given session is invalid (maybe timed out)
	ErrSessionInvalid = errors.New(`invalid session`)
)

type Store interface {
	// CreateSession returns a session token for the given user
	CreateSession(userId string) (string, error)

	// UserForSession returns the user associated with the given session, an error is returned if the session was missing or invalid
	UserForSession(sessionId string) (string, error)

	// DeleteSession removes the given session, an error is returned if the session was missing or invalid
	DeleteSession(sessionId string) error
}
