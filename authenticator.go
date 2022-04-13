package authsession

// Authenticator verifies user credentials and returns a list of permissions a user can have
type CredentialChecker interface {
	CheckCredentials(userId string, password string) (bool, error)
}

type PermissionChecker interface {
	HasPermissions(userId string, required []string) (bool, error)
}

type CredentialPermissionChecker interface {
	CredentialChecker
	PermissionChecker
}
