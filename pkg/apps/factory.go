package apps

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/api"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/jizzberry"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"net/http"
)

type App interface {
	Register(r *mux.Router)
}

var apps = make([]App, 0)

func RegisterApps(r *mux.Router) {

	apps = append(apps, api.Api{}, authentication.Authentication{}, jizzberry.Jizzberry{})

	for _, i := range apps {
		i.Register(r)
	}
}

func RegisterFileServer(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(pkger.Dir("/web/templates/static")),
	))

	r.PathPrefix("/thumbnails/").Handler(http.StripPrefix("/thumbnails/",
		http.FileServer(http.Dir(helpers.ThumbnailPath)),
	))
}
