package authsession

import "time"

func NewDefaultConfig() *CookieConfig {
	return &CookieConfig{
		Name:     "sid",
		Path:     "/",
		Duration: 7 * 24 * time.Hour,
	}
}

type CookieConfig struct {
	Name     string
	Path     string
	Duration time.Duration
}
