package jizzberry

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/auth"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/tags"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Jizzberry struct {
}

const component = "Web"

type Context struct {
	Files     []files.Files
	Tags      []tags.Tag
	Actors    []actor_details.ActorDetails
	ActorList string
	UpNext    []files.Files
	Config    helpers.Config
	Users     []auth.Auth
}

const baseURL = "/Jizzberry"

func (a Jizzberry) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()
	authRouter.HandleFunc("/home", homeHandler)
	authRouter.HandleFunc("/tags", allCategoriesHandler)
	authRouter.HandleFunc("/actors", allActorsHandler)
	authRouter.HandleFunc("/scene/{scene_id}", singleSceneHandler)
	authRouter.HandleFunc("/stream/{scene_id}", streamHandler)
	authRouter.HandleFunc("/settings", settingsHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allFiles := files.Initialize().Get(files.Files{})

	err := helpers.Render(w, http.StatusOK, "home", Context{Files: allFiles})
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allTags := tags.Initialize().Get(tags.Tag{})

	err := helpers.Render(w, http.StatusOK, "tags", Context{Tags: allTags})
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allActorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allActors := actor_details.Initialize().Get(actor_details.ActorDetails{})

	err := helpers.Render(w, http.StatusOK, "actors", Context{Actors: allActors})
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func singleSceneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	vars := mux.Vars(r)
	sceneId, _ := strconv.ParseInt(vars["scene_id"], 10, 64)

	file := files.Initialize().Get(files.Files{GeneratedID: sceneId})

	randomNext := files.Initialize().Get(files.Files{})
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(randomNext), func(i, j int) { randomNext[i], randomNext[j] = randomNext[j], randomNext[i] })

	// TODO: Get UpNext on same conn
	err := helpers.Render(w, http.StatusOK, "single", Context{
		Files:     file,
		ActorList: file[0].Actors,
		UpNext: func() []files.Files {
			if len(randomNext) > 9 {
				return randomNext[0:9]
			} else {
				return randomNext
			}
		}(),
	})

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allTags := tags.Initialize().Get(tags.Tag{})
	users := auth.Initialize().Get(auth.Auth{})
	err := helpers.Render(w, http.StatusOK, "settings", Context{Config: helpers.GetConfig(), Tags: allTags, Users: users})

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}
