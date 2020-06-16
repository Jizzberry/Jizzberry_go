package jizzberry

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/middleware"
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
	IsAdmin   bool
}

const baseURL = "/Jizzberry"

func (a Jizzberry) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	htmlRouter.Use(middleware.AuthMiddleware())

	htmlRouter.HandleFunc("/home", homeHandler)
	htmlRouter.HandleFunc("/tags", allCategoriesHandler)
	htmlRouter.HandleFunc("/actors", allActorsHandler)
	htmlRouter.HandleFunc("/actors/{actor_id}", singleActorHanlder)
	htmlRouter.HandleFunc("/scene/{scene_id}", singleSceneHandler)
	htmlRouter.HandleFunc("/stream/{scene_id}", streamHandler)
	htmlRouter.HandleFunc("/settings", settingsHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allFiles := files.Initialize().Get(files.Files{})

	ctx := Context{Files: allFiles}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "home", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allTags := tags.Initialize().Get(tags.Tag{})

	ctx := Context{Tags: allTags}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "tags", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allActorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	allActors := actor_details.Initialize().Get(actor_details.ActorDetails{})

	ctx := Context{Actors: allActors}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "actors", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func singleActorHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	vars := mux.Vars(r)
	actorIDstr, _ := vars["actor_id"]

	filesIDs := files.GetActorRelations(actorIDstr)

	filesModel := files.Initialize()

	ctx := Context{}
	sidebarContext(&ctx, r)

	for _, f := range filesIDs {
		i, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		ctx.Files = append(ctx.Files, filesModel.Get(files.Files{GeneratedID: i})...)
	}

	actorID, err := strconv.ParseInt(actorIDstr, 10, 64)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	actorDetails := actor_details.Initialize().Get(actor_details.ActorDetails{ActorId: actorID})
	ctx.Actors = actorDetails

	err = helpers.Render(w, http.StatusOK, "singleActor", ctx)
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
	ctx := Context{
		Files:     file,
		ActorList: file[0].Actors,
		UpNext: func() []files.Files {
			if len(randomNext) > 9 {
				return randomNext[0:9]
			} else {
				return randomNext
			}
		}(),
	}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "singleScene", ctx)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	if authentication.IsAdmin(authentication.GetUsernameFromSession(r)) {
		allTags := tags.Initialize().Get(tags.Tag{})
		users := auth.Initialize().Get(auth.Auth{})
		err := helpers.Render(w, http.StatusOK, "settings", Context{Config: helpers.GetConfig(), Tags: allTags, Users: users, IsAdmin: true})

		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
}

func sidebarContext(ctx *Context, r *http.Request) {
	if authentication.IsAdmin(authentication.GetUsernameFromSession(r)) {
		ctx.IsAdmin = true
	} else {
		ctx.IsAdmin = false
	}
}
