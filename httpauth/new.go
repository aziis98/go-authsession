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

func WithNotAuthorizedHandler(notAuthorizedHandler http.Handler) HttpOption {
	return HttpOptionFunc(func(ha *HttpAuth) {
		ha.NotAuthorizedHandler = notAuthorizedHandler
	})
}

func WithErrorHandler(errorHandler func(error) http.Handler) HttpOption {
	return HttpOptionFunc(func(ha *HttpAuth) {
		ha.ErrorHandler = errorHandler
	})
}

func defaultErrorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})
}

func New(credChecker authsession.CredentialChecker, opts ...authsession.Option) *HttpAuth {
	httpAuth := &HttpAuth{}
	httpAuth.ErrorHandler = defaultErrorHandler
	httpAuth.NotAuthorizedHandler = defaultErrorHandler(authsession.ErrNotAuthorized)

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
