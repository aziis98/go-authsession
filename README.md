# AuthSession

A library to easily handle cookie sessions and a basic form of permission management with go http frameworks and libraries.

For now this supports

-   [`net/http`](https://pkg.go.dev/net/http) using the submodule [`httpauth`](./httpauth).

-   [`fiber`](https://github.com/gofiber/fiber) using the submodule [`fiberauth`](./fiberauth).

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

We can now use directly `authsession.New(CredentialChecker, ...Option)` to create an instance of `*authsession.Base` that provides the following methods

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

Otherwise this can take some optional arguments for setting up error handlers used by the _middleware methods_ as shown later.

-   `WithUnauthorizedHandler(unauthorizedHandler http.Handler) Option`

    Special HTTP error handler for unauthorized access.

-   `WithErrorHandler(errorHandler func(error) http.Handler) Option`

    Generic HTTP error handler.

For example to login a user you can just define an http handler as follows

```go
http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})
```

Lastly if you are using a `net/http` based router you can use one of the following middlewares

```go
// by having the route accept only logged in users...
router.Use(auth.LoggedMiddleware())

// ...or by having the route require a list of permissions
router.Use(auth.PermissionsMiddleware([]string{ "moderator" }))
```

(for a complete example see [httpauth/auth_test.go](./httpauth/auth_test.go))

### Fiber Auth

The submodule [`fiberauth`](./fiberauth) provides an adapter working with the [`fiber`](https://github.com/gofiber/fiber) web framework. Let's use the `exampleAuth` struct from before.

TODO: Add an example

(for now see [fiberauth/auth.go](./fiberauth/auth.go))

## TODOs

-   [ ] Add token "refresh"-ability to `session.Store`.

-   [ ] Add some tests to `fiberauth`.
