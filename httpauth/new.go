package httpauth

import (
	"net/http"

	"github.com/aziis98/go-authsession"
)

type HttpOption interface {
	authsession.Option
	SetHttpOption(h *HttpAuth)
}

type HttpOptionFunc func(*HttpAuth)

func (_ HttpOptionFunc) SetOption(b *authsession.Base) {}

func (f HttpOptionFunc) SetHttpOption(h *HttpAuth) {
	f(h)
}

func WithUnauthorizedHandler(unauthorizedHandler http.Handler) authsession.Option {
	return HttpOptionFunc(func(ha *HttpAuth) {
		ha.unauthorizedHandler = unauthorizedHandler
	})
}

func WithErrorHandler(errorHandler func(error) http.Handler) authsession.Option {
	return HttpOptionFunc(func(ha *HttpAuth) {
		ha.errorHandler = errorHandler
	})
}

func defaultErrorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})
}

func New(credChecker authsession.CredentialChecker, opts ...authsession.Option) *HttpAuth {
	httpAuth := &HttpAuth{}
	httpAuth.errorHandler = defaultErrorHandler
	httpAuth.unauthorizedHandler = defaultErrorHandler(authsession.ErrUnauthorized)

	basicOpts := []authsession.Option{}
	for _, opt := range opts {
		if httpOpt, ok := opt.(HttpOption); ok {
			httpOpt.SetHttpOption(httpAuth)
		} else {
			basicOpts = append(basicOpts, opt)
		}
	}

	httpAuth.basic = authsession.New(credChecker, basicOpts...)

	return httpAuth
}
