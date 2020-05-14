package authentication

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Authentication struct {
}

type Login struct {
	Error string
}

const baseURL = "/auth"

var SessionsStore = sessions.NewCookieStore(config.GetSessionsKey())

func (a Authentication) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()

	authRouter.HandleFunc("/login/", loginHandler)
	authRouter.HandleFunc("/logout/", logoutHandler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	session, _ := SessionsStore.Get(r, config.SessionsKey)

	// If user is already logged in, don't show login page again until logout
	if ValidateSession(w, r) {
		http.Redirect(w, r, "/Jizzberry/", http.StatusFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	username := r.FormValue(config.Usernamekey)
	password := r.FormValue(config.PasswordKey)

	if username != "" && password != "" {
		if userIsValid(username, password) {

			session.Values[config.Usernamekey] = username
			prevURL := session.Values[config.PrevURLKey]

			session.Options.MaxAge = 30 * 60

			if prevURL != nil {
				session.Values[config.PrevURLKey] = nil
				err := session.Save(r, w)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, prevURL.(string), http.StatusFound)
				return
			}

			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/Jizzberry/", http.StatusFound)
			return
		}

		err := helpers.Render(w, http.StatusOK, "login", Login{Error: "Couldn't validate"})
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	err = helpers.Render(w, http.StatusOK, "login", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := SessionsStore.Get(r, config.SessionsKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, config.Usernamekey)
	session.Options.MaxAge = -1

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, config.LoginURL, http.StatusFound)
}

func userIsValid(username string, password string) bool {
	fetchUsers := auth.Initialize().Get(auth.Auth{Username: username})
	if len(fetchUsers) > 0 {
		hashedPass := fetchUsers[0].Password
		if hashedPass != "" {
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
			if err != nil {
				fmt.Println(err)
				return false
			}
			return true
		}
	}
	return false
}

func ValidateSession(w http.ResponseWriter, r *http.Request) bool {
	session, err := SessionsStore.Get(r, config.SessionsKey)
	if err != nil {
		return false
	}

	if session.IsNew {
		return false
	}

	val := session.Values[config.Usernamekey]

	if val != nil {
		user := auth.Initialize().Get(auth.Auth{Username: val.(string)})
		if len(user) > 0 {
			if user[0].Username == val {
				session.Options.MaxAge = 30 * 60
				err := session.Save(r, w)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return true
			}
		}
	}
	return false
}
