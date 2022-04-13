# AuthSession

## Usage

First we need something implementing the `authsession.CredentialChecker` interface.

```go
type exampleAuth struct{}

func (_ *exampleAuth) CheckCredentials(userId string, password string) (bool, error) {
	if userId != "example" {
        // return authsession.ErrUserNotFound or nil as you prefer
		return false, authsession.ErrUserNotFound
	}

	return password == "123", nil
}
```

We can now use directly `authsession.New(CredentialChecker, ...Option)` to create an instance of `*authsession.Base` that provides the following methods.

-   `(*Base).Login(userId string, password string) (string, error)`

    If the provided credentials are correct this returns a new session token for this user.

-   `(*Base).Logout(sessionId string) error`

    If the given session is valid this method deletes the given session.

-   `(*Base).UserForSession(sessionId string) (string, error)`

    Retrieves the user associated with this session.

-   `(*Base).IsLogged(sessionId string) (bool, error)`

    This uses the previous method to check if the user is logged given a session token.

-   `(*Base).HasPermissions(userId string, required []string) (bool, error)`

    This function can only be used if the instance of `CredentialChecker` is also a `PermissionChecker` and panics otherwise.

(A valid user and a valid session token must be a non empty string)

Otherwise instead of using `*authsession.Base` directly you can use one of the following adapters for various libraries.

### Http Auth

The submodule [`httpauth`](./httpauth) provides an adapter working with the `net/http` module. Let's use the `exampleAuth` struct from before.

The http adapter can be created using `httpauth.New(CredentialChecker, ...Option)`. The simplest form to initialize it is the following

```go
auth := httpauth.New(&exampleAuth{})
```

Otherwise this can take two optional arguments for setting up error handlers used by the middleware methods.

-   `httpauth.WithNotAuthorizedHandler(notAuthorizedHandler http.Handler) Option`

    Special HTTP error handler for unauthorized access.

-   `httpauth.WithErrorHandler(errorHandler func(error) http.Handler) Option`

    Generic HTTP error handler.

For example to login in a user you can just define an http handler as follows

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    if err := auth.Login(w, username, password); err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    if err := json.NewEncoder(w).Encode("ok"); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}
```

Lastly if you are using a `net/http` based router you can use one of the following middlewares

```go
router.Use(auth.LoggedMiddleware())

// or...
router.Use(auth.PermissionsMiddleware([]string{ "moderator" }))
```

### Fiber Auth

The submodule [`fiberauth`](./fiberauth) provides an adapter working with the [`fiber`](https://github.com/gofiber/fiber) web framework. Let's use the `exampleAuth` struct from before.

...

(for now see [fiberauth/auth.go](./fiberauth/auth.go))

## TODOs

-   [ ] Add token "refresh"-ability to `session.Store`.

-   [ ] Add a test to `fiberauth`.
