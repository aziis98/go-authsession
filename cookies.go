package authsession

import "time"

type CookieAdapter interface {
	Cookie(ctx any, name string) (string, error)
	SetCookie(ctx any, cookie Cookie) error
}

type Cookie struct {
	Name    string
	Path    string
	Value   string
	Expires time.Time
}
