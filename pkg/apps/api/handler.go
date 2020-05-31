package api

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/middleware"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	studios2 "github.com/Jizzberry/Jizzberry-go/pkg/models/studios"
	tags2 "github.com/Jizzberry/Jizzberry-go/pkg/models/tags"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/manager"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Api struct {
}

type task struct {
	Uid string `json:"uid"`
}

type progress struct {
	Progress int `json:"progress"`
}

func (a Api) Register(r *mux.Router) {

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.StrictSlash(false)

	apiRouter.Use(middleware.AuthMiddleware())

	apiRouter.HandleFunc("/files", filesHandler).Methods("GET")
	apiRouter.HandleFunc("/actor_details", actorDetailHandler).Methods("GET")
	apiRouter.HandleFunc("/actors", actorsHandler).Methods("GET")
	apiRouter.HandleFunc("/studios", studiosHandler).Methods("GET")
	apiRouter.HandleFunc("/tags", tagsHandler).Methods("GET")
	apiRouter.HandleFunc("/scrapeActors", scrapeActorHandler).Methods("GET")
	apiRouter.HandleFunc("/startScanTask", scanHandler).Methods("POST", "GET")
	apiRouter.HandleFunc("/startScrapeTask", scrapeListHandler).Methods("POST", "GET")
	apiRouter.HandleFunc("/progress", getProgress).Methods("GET")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	var file []files.Files
	model := files.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		file = model.Get(files.Files{GeneratedID: int64(genId)})

	} else if len(queryParams["file_name"]) > 0 {
		file = model.Get(files.Files{FileName: "%" + queryParams["file_name"][0] + "%"})

	} else if len(queryParams["file_path"]) > 0 {
		file = model.Get(files.Files{FileName: "%" + queryParams["file_path"][0] + "%"})

	} else {
		file = model.Get(files.Files{})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&file)
	if err != nil {
		fmt.Println(err)
	}
}

func actorDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	var actorDetails []actor_details.ActorDetails
	model := actor_details.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		actorDetails = model.Get(actor_details.ActorDetails{GeneratedId: int64(genId)})

	} else if len(queryParams["name"]) > 0 {
		actorDetails = model.Get(actor_details.ActorDetails{Name: "%" + queryParams["name"][0] + "%"})

	} else if len(queryParams["scene_id"]) > 0 {
		sceneId, err := strconv.Atoi(queryParams["scene_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		actorDetails = model.Get(actor_details.ActorDetails{SceneId: int64(sceneId)})

	} else if len(queryParams["actor_id"]) > 0 {
		actorId, err := strconv.Atoi(queryParams["actor_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		actorDetails = model.Get(actor_details.ActorDetails{SceneId: int64(actorId)})

	} else {
		actorDetails = model.Get(actor_details.ActorDetails{})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&actorDetails)
	if err != nil {
		fmt.Println(err)
	}
}

func actorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	actors := make([]actor.Actor, 0)
	model := actor.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		actors = model.Get(actor.Actor{GeneratedID: int64(genId)})
	} else if len(queryParams["name"]) > 0 {
		actors = model.Get(actor.Actor{Name: "%" + queryParams["name"][0] + "%"})
	} else if len(queryParams["title"]) > 0 {
		actors = tasks.MatchActorToTitle(queryParams["title"][0])
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&actors)
	if err != nil {
		fmt.Println(err)
	}
}

func scrapeActorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	actors := make([]actor_details.ActorDetails, 0)

	if len(queryParams["actor_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["actor_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		tmp := actor.Initialize().Get(actor.Actor{GeneratedID: int64(genId)})
		if len(tmp) > 0 {
			actors = append(actors, *scrapers.ScrapeActor(0, tmp[0]))
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&actors)
	if err != nil {
		fmt.Println(err)
	}
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := manager.StartScan()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	encoder.Encode(&task{Uid: uid})
}

func scrapeListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := make([]task, 0)
	t = append(t, task{Uid: manager.StartScrapeActors()})
	t = append(t, task{Uid: manager.StartScrapeStudios()})
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	encoder.Encode(&t)
}

func studiosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	studios := make([]studios2.Studios, 0)
	model := studios2.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		studios = model.Get(studios2.Studios{GeneratedID: int64(genId)})
	} else if len(queryParams["name"]) > 0 {
		studios = model.Get(studios2.Studios{Studio: "%" + queryParams["name"][0] + "%"})
	} else {
		studios = model.Get(studios2.Studios{})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&studios)
	if err != nil {
		fmt.Println(err)
	}
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	tags := make([]tags2.Tags, 0)
	model := tags2.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		tags = model.Get(tags2.Tags{GeneratedID: int64(genId)})
	} else if len(queryParams["name"]) > 0 {
		tags = model.Get(tags2.Tags{Name: "%" + queryParams["name"][0] + "%"})
	} else {
		tags = model.Get(tags2.Tags{})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&tags)
	if err != nil {
		fmt.Println(err)
	}
}

func getProgress(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	progress := progress{}
	progress.Progress = -1
	if len(queryParams["uid"]) > 0 {
		progress.Progress = manager.GetProgress(queryParams["uid"][0])
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&progress)
	if err != nil {
		fmt.Println(err)
	}
}
