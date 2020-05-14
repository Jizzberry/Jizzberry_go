package middleware

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authentication.ValidateSession(w, r) {
				next.ServeHTTP(w, r)
			} else {
				session, _ := authentication.SessionsStore.Get(r, config.SessionsKey)
				session.Values[config.PrevURLKey] = r.URL.Path
				err := session.Save(r, w)
				if err != nil {
					fmt.Println(err)
				}
				http.Redirect(w, r, "/auth/login/", http.StatusFound)
			}
		})
	}
}
