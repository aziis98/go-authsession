package authsession

import (
	"time"

	"github.com/aziis98/go-authsession/session"
)

type OptionFunc func(*Base)

func (opt OptionFunc) SetOption(as *Base) {
	opt(as)
}

type Option interface {
	SetOption(as *Base)
}

func WithConfig(name, path string, duration time.Duration) Option {
	return OptionFunc(func(as *Base) {
		as.Config = &CookieConfig{name, path, duration}
	})
}

func WithSessionStore(store session.Store) Option {
	return OptionFunc(func(as *Base) {
		as.sessionStore = store
	})
}

func New(credChecker CredentialChecker, opts ...Option) *Base {
	authSession := &Base{
		Config:      NewDefaultConfig(),
		credChecker: credChecker,
	}

	if permChecker, ok := credChecker.(PermissionChecker); ok {
		authSession.permChecker = permChecker
	}

	for _, opt := range opts {
		opt.SetOption(authSession)
	}

	// by deafult if no session store is provided use an in memory map
	if authSession.sessionStore == nil {
		authSession.sessionStore = session.NewInMemoryStore()
	}

	return authSession
}
