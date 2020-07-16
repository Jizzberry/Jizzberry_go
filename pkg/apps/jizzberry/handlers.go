package jizzberry

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/apps/authentication"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/middleware"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/auth"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/studios"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/tags"
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
	Studios   []studios.Studio
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
	htmlRouter.HandleFunc("/tags/{tag_id}", singleTagHandler)
	htmlRouter.HandleFunc("/actors", allActorsHandler)
	htmlRouter.HandleFunc("/actors/{actor_id}", singleActorHanlder)
	htmlRouter.HandleFunc("/studios", allStudiosHandler)
	htmlRouter.HandleFunc("/studios/{studio_id}", singleStudiosHanlder)
	htmlRouter.HandleFunc("/scene/{scene_id}", singleSceneHandler)
	htmlRouter.HandleFunc("/stream/{scene_id}", streamHandler)
	htmlRouter.HandleFunc("/settings", settingsHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	model := files.Initialize()
	defer model.Close()

	allFiles := model.Get(files.Files{})

	ctx := Context{Files: allFiles}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "home", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	model := tags.Initialize()
	defer model.Close()

	allTags := model.Get(tags.Tag{})

	ctx := Context{Tags: allTags}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "tags", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func parseFilesFromRelation(fileIds []string, ctx *Context) {
	filesModel := files.Initialize()
	defer filesModel.Close()

	for _, f := range fileIds {
		i, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		ctx.Files = append(ctx.Files, filesModel.Get(files.Files{GeneratedID: i})...)
	}
}

func singleTagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	ctx := Context{}
	sidebarContext(&ctx, r)

	vars := mux.Vars(r)
	tagIdstr, _ := vars["tag_id"]

	filesIDs := files.GetTagRelations(tagIdstr)

	tagId, err := strconv.ParseInt(tagIdstr, 10, 64)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	model := tags.Initialize()
	defer model.Close()

	tagDetails := model.Get(tags.Tag{GeneratedID: tagId})
	ctx.Tags = tagDetails

	parseFilesFromRelation(filesIDs, &ctx)

	err = helpers.Render(w, http.StatusOK, "singleTag", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allActorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	model := actor_details.Initialize()
	defer model.Close()

	allActors := make([]actor_details.ActorDetails, 0)

	for _, k := range files.GetUsedActors() {
		key, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allActors = append(allActors, model.Get(actor_details.ActorDetails{ActorId: key})...)
	}

	ctx := Context{Actors: allActors}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "actors", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func singleActorHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	ctx := Context{}
	sidebarContext(&ctx, r)

	vars := mux.Vars(r)
	actorIDstr, _ := vars["actor_id"]

	filesIDs := files.GetActorRelations(actorIDstr)

	actorID, err := strconv.ParseInt(actorIDstr, 10, 64)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	model := actor_details.Initialize()
	defer model.Close()

	actorDetails := model.Get(actor_details.ActorDetails{ActorId: actorID})
	ctx.Actors = actorDetails

	parseFilesFromRelation(filesIDs, &ctx)

	err = helpers.Render(w, http.StatusOK, "singleActor", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func allStudiosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	model := studios.Initialize()
	defer model.Close()

	allStudios := make([]studios.Studio, 0)

	keys := files.GetUsedStudios()
	for _, k := range keys {
		key, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allStudios = append(allStudios, model.Get(studios.Studio{GeneratedID: key})...)
	}

	ctx := Context{Studios: allStudios}
	sidebarContext(&ctx, r)

	err := helpers.Render(w, http.StatusOK, "studios", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func singleStudiosHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	ctx := Context{}
	sidebarContext(&ctx, r)

	vars := mux.Vars(r)
	studioIDstr, _ := vars["studio_id"]

	filesIDs := files.GetStudioRelations(studioIDstr)

	studioID, err := strconv.ParseInt(studioIDstr, 10, 64)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	model := studios.Initialize()
	defer model.Close()

	studioDetails := model.Get(studios.Studio{GeneratedID: studioID})
	ctx.Studios = studioDetails

	parseFilesFromRelation(filesIDs, &ctx)

	err = helpers.Render(w, http.StatusOK, "singleStudios", ctx)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func singleSceneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	vars := mux.Vars(r)
	sceneId, _ := strconv.ParseInt(vars["scene_id"], 10, 64)

	model := files.Initialize()
	defer model.Close()

	file := model.Get(files.Files{GeneratedID: sceneId})
	if len(file) <= 0 {
		return
	}

	randomNext := model.Get(files.Files{})
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(randomNext), func(i, j int) { randomNext[i], randomNext[j] = randomNext[j], randomNext[i] })

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

	if authentication.IsAdminFromSession(r) {
		modelTags := tags.Initialize()
		defer modelTags.Close()

		allTags := modelTags.Get(tags.Tag{})

		modelAuth := auth.Initialize()
		defer modelAuth.Close()

		users := modelAuth.Get(auth.Auth{})

		err := helpers.Render(w, http.StatusOK, "settings", Context{Config: helpers.GetConfig(), Tags: allTags, Users: users, IsAdmin: true})

		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
}

func sidebarContext(ctx *Context, r *http.Request) {
	if authentication.IsAdminFromSession(r) {
		ctx.IsAdmin = true
	} else {
		ctx.IsAdmin = false
	}
}
