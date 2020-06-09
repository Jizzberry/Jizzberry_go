package authentication

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Authentication struct {
}

type Context struct {
	Error string
}

const baseURL = "/auth"

var SessionsStore = sessions.NewCookieStore(helpers.GetSessionsKey())

func (a Authentication) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()

	authRouter.HandleFunc("/login/", loginHandler)
	authRouter.HandleFunc("/logout/", logoutHandler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	session, _ := SessionsStore.Get(r, helpers.SessionsKey)

	// If user is already logged in, don't show login page again until logout
	if ValidateSession(w, r) {
		http.Redirect(w, r, "/Jizzberry/", http.StatusFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	username := r.FormValue(helpers.Usernamekey)
	password := r.FormValue(helpers.PasswordKey)

	if username != "" && password != "" {
		if userIsValid(username, password) {

			fmt.Println("valid")

			session.Values[helpers.Usernamekey] = username
			prevURL := session.Values[helpers.PrevURLKey]

			session.Options.MaxAge = 30 * 60

			if prevURL != nil {
				session.Values[helpers.PrevURLKey] = nil
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

		err := helpers.Render(w, http.StatusOK, "login", Context{Error: "Couldn't validate"})
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
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, helpers.Usernamekey)
	session.Options.MaxAge = -1

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, helpers.LoginURL, http.StatusFound)
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
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		return false
	}

	if session.IsNew {
		return false
	}

	val := session.Values[helpers.Usernamekey]

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
