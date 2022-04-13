package authsession

import (
	"github.com/aziis98/go-authsession/generics"
	"github.com/google/uuid"
)

type SessionStore[UserId comparable] interface {
	CreateSession(userId UserId) (string, error)
	UserForSession(sid string) (UserId, error)
	DeleteSession(sid string) error
}

type InMemoryStore[UserId comparable] map[string]UserId

var _ SessionStore[string] = InMemoryStore[string]{}

func NewInMemoryStore[UserId comparable]() InMemoryStore[UserId] {
	return make(InMemoryStore[UserId])
}

func (mem InMemoryStore[UserId]) CreateSession(userId UserId) (string, error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	sid := newId.String()
	mem[sid] = userId

	return sid, nil
}

func (mem InMemoryStore[UserId]) DeleteSession(sid string) error {
	if _, present := mem[sid]; !present {
		return ErrInvalidSession
	}

	delete(mem, sid)

	return nil
}

func (mem InMemoryStore[UserId]) UserForSession(sid string) (UserId, error) {
	userId, present := mem[sid]
	if !present {
		return generics.Zero[UserId](), ErrInvalidSession
	}

	return userId, nil
}
