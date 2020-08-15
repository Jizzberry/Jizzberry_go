package middleware

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

// Validates user from session
func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authentication.ValidateSession(w, r) {
				// If session is valid then pass to next handler
				next.ServeHTTP(w, r)
			} else {
				// If session is invalid
				session, _ := authentication.SessionsStore.Get(r, helpers.SessionsKey)
				// Store url in session to redirect after login
				session.Values[helpers.PrevURLKey] = r.URL.Path
				err := session.Save(r, w)
				if err != nil {
					helpers.LogError(err.Error())
				}
				http.Redirect(w, r, helpers.LoginURL, http.StatusFound)
			}
		})
	}
}
