package middleware

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

const component = "Middleware"

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authentication.ValidateSession(w, r) {
				next.ServeHTTP(w, r)
			} else {
				session, _ := authentication.SessionsStore.Get(r, helpers.SessionsKey)
				session.Values[helpers.PrevURLKey] = r.URL.Path
				err := session.Save(r, w)
				if err != nil {
					helpers.LogError(err.Error(), component)
				}
				http.Redirect(w, r, "/auth/login/", http.StatusFound)
			}
		})
	}
}
