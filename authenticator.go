package authsession

// CredentialChecker checks credentials given by a user
type CredentialChecker interface {
	CheckCredentials(userId string, password string) (bool, error)
}

// PermissionChecker checks if a user has some permissions
type PermissionChecker interface {
	HasPermissions(userId string, required []string) (bool, error)
}

type CredentialPermissionChecker interface {
	CredentialChecker
	PermissionChecker
}
