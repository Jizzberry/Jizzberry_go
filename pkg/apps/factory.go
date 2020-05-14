package apps

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/api"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/authentication"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"net/http"
)

type App interface {
	Register(r *mux.Router)
}

var apps = make([]App, 0)

func RegisterApps(r *mux.Router) {

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(pkger.Dir("/web/templates/static")),
	))

	apps = append(apps, api.Api{}, authentication.Authentication{})

	for _, i := range apps {
		i.Register(r)
	}
}
