package session

type Store interface {
	// CreateSession returns a session token for the given user
	CreateSession(userId string) (string, error)

	// UserForSession returns the user associated with the given session, an error is returned if the session was missing or invalid
	UserForSession(sessionId string) (string, error)

	// DeleteSession removes the given session, an error is returned if the session was missing or invalid
	DeleteSession(sessionId string) error
}
