package httpadapter

import (
	"fmt"
	"net/http"

	"github.com/aziis98/go-authsession"
)

type httpdriver struct{}

func (_ *httpdriver) Cookie(ctx any, name string) (string, error) {
	r, ok := ctx.(*http.Request)
	if !ok {
		return "", fmt.Errorf(`invalid context, must be of type "*http.Request"`)
	}

	c, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

func (_ *httpdriver) SetCookie(ctx any, cookie authsession.Cookie) error {
	w, ok := ctx.(http.ResponseWriter)
	if !ok {
		return fmt.Errorf(`invalid context, must be of type "*http.Request"`)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    cookie.Name,
		Path:    cookie.Path,
		Value:   cookie.Value,
		Expires: cookie.Expires,
	})

	return nil
}

func New() authsession.CookieAdapter {
	return &httpdriver{}
}
