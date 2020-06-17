package authentication

import (
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

const (
	baseURL   = "/auth"
	component = "WebAuth"
)

var SessionsStore = sessions.NewCookieStore(helpers.GetSessionsKey())

func (a Authentication) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()

	authRouter.HandleFunc("/login/", loginHandler)
	authRouter.HandleFunc("/logout/", logoutHandler)
	authRouter.HandleFunc("/create/", newUser).Methods("POST")
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
		helpers.LogError(err.Error(), component)
	}

	username := r.FormValue(helpers.Usernamekey)
	password := r.FormValue(helpers.PasswordKey)

	if username != "" && password != "" {
		if userIsValid(username, password) {
			session.Values[helpers.Usernamekey] = username
			prevURL := session.Values[helpers.PrevURLKey]

			session.Options.MaxAge = 30 * 60

			if prevURL != nil {
				session.Values[helpers.PrevURLKey] = nil
				err := session.Save(r, w)

				if err != nil {
					helpers.LogError(err.Error(), component)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, prevURL.(string), http.StatusFound)
				return
			}

			err := session.Save(r, w)
			if err != nil {
				helpers.LogError(err.Error(), component)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/Jizzberry/", http.StatusFound)
			return
		}

		err := helpers.Render(w, http.StatusOK, "login", Context{Error: "Couldn't validate"})
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return
	}
	err = helpers.Render(w, http.StatusOK, "login", nil)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		helpers.LogError(err.Error(), component)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, helpers.Usernamekey)
	session.Options.MaxAge = -1

	err = session.Save(r, w)

	if err != nil {
		helpers.LogError(err.Error(), component)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, helpers.LoginURL, http.StatusFound)
}

func userIsValid(username string, password string) bool {
	model := auth.Initialize()
	defer model.Close()

	fetchUsers := model.Get(auth.Auth{Username: username})
	if len(fetchUsers) > 0 {
		hashedPass := fetchUsers[0].Password
		if hashedPass != "" {
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
			if err != nil {
				helpers.LogError(err.Error(), component)
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
		model := auth.Initialize()
		defer model.Close()

		user := model.Get(auth.Auth{Username: val.(string)})
		if len(user) > 0 {
			if user[0].Username == val {
				session.Options.MaxAge = 30 * 60
				err := session.Save(r, w)
				if err != nil {
					helpers.LogError(err.Error(), component)
					return false
				}
				return true
			}
		}
	}
	return false
}

func newUser(w http.ResponseWriter, r *http.Request) {
	if !ValidateSession(w, r) && IsAdminFromSession(r) {
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	username := r.FormValue(helpers.Usernamekey)
	password := r.FormValue(helpers.PasswordKey)
	admin := r.FormValue("isAdmin")

	model := auth.Initialize()
	defer model.Close()

	model.Create(auth.Auth{
		Username: username,
		Password: password,
		IsAdmin:  admin == "on",
	})
}

func GetUsernameFromSession(r *http.Request) string {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		return ""
	}
	if session.IsNew {
		return ""
	}

	return session.Values[helpers.Usernamekey].(string)
}

func IsAdminFromSession(r *http.Request) bool {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		return false
	}
	if session.IsNew {
		return false
	}
	model := auth.Initialize()
	defer model.Close()

	user := model.Get(auth.Auth{Username: session.Values[helpers.Usernamekey].(string)})
	if len(user) > 0 {
		if user[0].IsAdmin {
			return true
		}
	}
	return false
}
