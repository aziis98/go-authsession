package fiberauth

import "github.com/aziis98/go-authsession"

func New(credChecker authsession.CredentialChecker, opts ...authsession.Option) *FiberAuth {
	fiberAuth := &FiberAuth{}
	fiberAuth.basic = authsession.New(credChecker, opts...)

	return fiberAuth
}
