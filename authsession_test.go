package authsession_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aziis98/go-authsession"
	"github.com/aziis98/go-authsession/httpadapter"
)

func setupAuthSessionService() *authsession.Service[string] {
	return &authsession.Service[string]{
		authsession.DefaultConfig,
		&authsession.AuthenticatorFunc[string]{
			func(username string) (string, error) {
				return username, nil
			},
			func(username, password string) (bool, error) {
				if username != "example" {
					return false, nil
				}

				return true, nil
			},
			func(id string, required []string) (bool, error) {
				if id == "example" {
					return true, nil
				}

				return false, authsession.ErrNotAuthorized
			},
		},
		authsession.NewInMemoryStore[string](),
		httpadapter.New(),
	}
}

func TestLogin(t *testing.T) {
	auth := setupAuthSessionService()

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

	log.Print(res.Result().Cookies())
}

func TestLogout(t *testing.T) {
	auth := setupAuthSessionService()

	logoutHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := auth.Logout(w); err != nil {
			t.Fatal(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	res := httptest.NewRecorder()

	logoutHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatal()
	}

	log.Printf(`%+v`, res.Result())
}
