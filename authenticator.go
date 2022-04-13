package authsession

// Authenticator verifies user credentials and returns a list of permissions a user can have
type Authenticator[UserId comparable] interface {
	FromUsername(username string) (UserId, error)
	Authenticate(username, password string) (bool, error)
	HasPermissions(id UserId, required []string) (bool, error)
}

type AuthenticatorFunc[UserId comparable] struct {
	FromUsernameFn   func(username string) (UserId, error)
	AuthenticateFn   func(username, password string) (bool, error)
	HasPermissionsFn func(id UserId, required []string) (bool, error)
}

func (af AuthenticatorFunc[UserId]) FromUsername(username string) (UserId, error) {
	return af.FromUsernameFn(username)
}

func (af AuthenticatorFunc[UserId]) Authenticate(username, password string) (bool, error) {
	return af.AuthenticateFn(username, password)
}

func (af AuthenticatorFunc[UserId]) HasPermissions(id UserId, required []string) (bool, error) {
	return af.HasPermissionsFn(id, required)
}
