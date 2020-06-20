package apps

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/api"
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/jizzberry"
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/websocket"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
)

type App interface {
	Register(r *mux.Router)
}

var apps = []App{api.Api{}, authentication.Authentication{}, jizzberry.Jizzberry{}, websocket.Websocket{}}

func RegisterApps(r *mux.Router) {
	for _, i := range apps {
		i.Register(r)
	}
}

func RegisterFileServer(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(filepath.Join(helpers.GetWorkingDirectory(), "web/templates/static"))),
	))

	r.PathPrefix("/thumbnails/").Handler(http.StripPrefix("/thumbnails/",
		http.FileServer(http.Dir(helpers.ThumbnailPath)),
	))

	r.PathPrefix("/logs/").Handler(http.StripPrefix("/logs/",
		http.FileServer(http.Dir(helpers.LogDir)),
	))
}
