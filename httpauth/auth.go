package httpauth

import (
	"net/http"
	"time"

	"github.com/aziis98/go-authsession"
)

type HttpAuth struct {
	basic *authsession.Base

	unauthorizedHandler http.Handler
	errorHandler        func(error) http.Handler
}

func (h *HttpAuth) cookie(r *http.Request) (string, error) {
	c, err := r.Cookie(h.basic.Config.Name)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

func (h *HttpAuth) setCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.basic.Config.Name,
		Path:    h.basic.Config.Path,
		Value:   value,
		Expires: time.Now().Add(h.basic.Config.Duration),
	})
}

func (h *HttpAuth) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.basic.Config.Name,
		Path:    h.basic.Config.Path,
		Value:   "",
		Expires: time.Now(),
	})
}

func (h *HttpAuth) Login(w http.ResponseWriter, username, password string) error {
	sid, err := h.basic.Login(username, password)
	if err != nil {
		return err
	}

	h.setCookie(w, sid)

	return nil
}

func (h *HttpAuth) Logout(w http.ResponseWriter, r *http.Request) error {
	sessionId, err := h.cookie(r)
	if err != nil {
		return err
	}

	if _, err := h.basic.UserForSession(sessionId); err != nil {
		return err
	}

	h.clearCookie(w)

	return nil
}

func (h *HttpAuth) IsLogged(r *http.Request) (bool, error) {
	sessionId, err := h.cookie(r)
	if err != nil {
		return false, err
	}

	if _, err := h.basic.UserForSession(sessionId); err != nil {
		return false, err
	}

	return true, nil
}

func (h *HttpAuth) LoggedMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logged, err := h.IsLogged(r)
			if err != nil {
				h.errorHandler(err).ServeHTTP(w, r)
				return
			}

			if !logged {
				h.unauthorizedHandler.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *HttpAuth) PermissionsMiddleware(required []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sid, err := h.cookie(r)
			if err != nil {
				h.errorHandler(err).ServeHTTP(w, r)
				return
			}

			userId, err := h.basic.UserForSession(sid)
			if err != nil {
				h.errorHandler(err).ServeHTTP(w, r)
				return
			}

			ok, err := h.basic.HasPermissions(userId, required)
			if err != nil {
				h.errorHandler(err).ServeHTTP(w, r)
				return
			}

			if !ok {
				h.unauthorizedHandler.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
