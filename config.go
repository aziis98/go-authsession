package authsession

import "time"

var DefaultConfig = ServiceConfig{
	CookieName:     "sid",
	CookiePath:     "/",
	CookieDuration: 7 * 24 * time.Hour,
}

type ServiceConfig struct {
	CookieName     string
	CookiePath     string
	CookieDuration time.Duration
}
