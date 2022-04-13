package httpauth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aziis98/go-authsession"
	"github.com/aziis98/go-authsession/httpauth"
)

type exampleAuth struct{}

func (_ *exampleAuth) CheckCredentials(userId string, password string) (bool, error) {
	if userId != "example" {
		return false, nil
	}

	return password == "123", nil
}

func (_ *exampleAuth) HasPermissions(id string, required []string) (bool, error) {
	if id == "example" {
		for _, perm := range required {
			if perm != "admin" {
				return false, nil
			}
		}

		return true, nil
	}

	return false, authsession.ErrNotAuthorized
}

var _ authsession.CredentialPermissionChecker = &exampleAuth{}

func TestLogin(t *testing.T) {
	auth := httpauth.New(&exampleAuth{})

	var cookie1 *http.Cookie
	t.Run("Login", func(t *testing.T) {
		loginHandler := func(w http.ResponseWriter, r *http.Request) {
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

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`username=example&password=123`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res := httptest.NewRecorder()

		loginHandler(res, req)

		if res.Code != http.StatusOK {
			t.Fatal()
		}

		cookie1 = res.Result().Cookies()[0]
		t.Log(res.Result().Cookies())
	})

	t.Run("Logout", func(t *testing.T) {
		logoutHandler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if err := auth.Logout(w, r); err != nil {
				t.Fatal(err)
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(cookie1)

		res := httptest.NewRecorder()

		logoutHandler(res, req)

		if res.Code != http.StatusOK {
			t.Fatal()
		}

		t.Log(res.Result().Cookies())
	})
}

// func TestLogout(t *testing.T) {
// 	auth := setupAuthSessionService()

// }
