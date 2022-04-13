package session

import "errors"

var (
	// ErrSessionNotFound is raised when the given session is not present in the session store
	ErrSessionNotFound = errors.New(`session not found`)

	// ErrSessionInvalid is raised when the given session is invalid (maybe timed out)
	ErrSessionInvalid = errors.New(`invalid session`)
)
