package session

import (
	"github.com/google/uuid"
)

// mapStore is a implemented as a simple map from session id to user id
type mapStore map[string]string

func NewInMemoryStore() Store {
	return make(mapStore)
}

func (mem mapStore) CreateSession(userId string) (string, error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	newSessionId := newId.String()
	mem[newSessionId] = userId

	return newSessionId, nil
}

func (mem mapStore) DeleteSession(sessionId string) error {
	if _, present := mem[sessionId]; !present {
		return ErrSessionNotFound
	}

	delete(mem, sessionId)

	return nil
}

func (mem mapStore) UserForSession(sessionId string) (string, error) {
	userId, present := mem[sessionId]
	if !present {
		return "", ErrSessionNotFound
	}

	return userId, nil
}
