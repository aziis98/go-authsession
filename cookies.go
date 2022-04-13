package authsession

import "time"

type Cookie struct {
	Name    string
	Path    string
	Value   string
	Expires time.Time
}
