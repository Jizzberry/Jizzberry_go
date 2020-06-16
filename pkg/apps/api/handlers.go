package api

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/middleware"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	studios2 "github.com/Jizzberry/Jizzberry-go/pkg/models/studios"
	tags2 "github.com/Jizzberry/Jizzberry-go/pkg/models/tags"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/manager"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const component = "API"

type Api struct {
}

type task struct {
	Uid string `json:"uid"`
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
	apiRouter.HandleFunc("/stopTask", stopHandler).Methods("POST")
	apiRouter.HandleFunc("/config", configHandler).Methods("GET", "POST")
	apiRouter.HandleFunc("/setPath", pathHandler).Methods("DELETE", "POST")
	apiRouter.HandleFunc("/metadata", parseMetadata).Methods("POST")
	apiRouter.HandleFunc("/browse", fileBrowser).Methods("GET")
	apiRouter.HandleFunc("/queryScrapers", queryScrapersHandler).Methods("GET")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	var file []files.Files
	model := files.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			helpers.LogError(err.Error(), component)
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
		helpers.LogError(err.Error(), component)
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
			helpers.LogError(err.Error(), component)
		}
		actorDetails = model.Get(actor_details.ActorDetails{GeneratedId: int64(genId)})

	} else if len(queryParams["name"]) > 0 {
		actorDetails = model.Get(actor_details.ActorDetails{Name: "%" + queryParams["name"][0] + "%"})

	} else if len(queryParams["actor_id"]) > 0 {
		actorId, err := strconv.Atoi(queryParams["actor_id"][0])
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		actorDetails = model.Get(actor_details.ActorDetails{ActorId: int64(actorId)})

	} else {
		actorDetails = model.Get(actor_details.ActorDetails{})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&actorDetails)
	if err != nil {
		helpers.LogError(err.Error(), component)
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
			helpers.LogError(err.Error(), component)
		}
		actors = model.Get(actor.Actor{GeneratedID: int64(genId)})
	} else if len(queryParams["name"]) > 0 {
		actors = model.Get(actor.Actor{Name: "%" + queryParams["name"][0] + "%"})
	} else if len(queryParams["title"]) > 0 {
		actors = tasks_handler.MatchActorToTitle(queryParams["title"][0])
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&actors)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func scrapeActorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	actors := make([]actor_details.ActorDetails, 0)

	if len(queryParams["actor_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["actor_id"][0])
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		tmp := actor.Initialize().Get(actor.Actor{GeneratedID: int64(genId)})
		if len(tmp) > 0 {
			actors = append(actors, *scrapers.ScrapeActor(tmp[0]))
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&actors)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func scanHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := manager.StartScan()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&task{Uid: uid})
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func scrapeListHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := make([]task, 0)
	t = append(t, task{Uid: manager.StartScrapeActors()})
	t = append(t, task{Uid: manager.StartScrapeStudios()})
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&t)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func studiosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	studios := make([]studios2.Studio, 0)
	model := studios2.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		studios = model.Get(studios2.Studio{GeneratedID: int64(genId)})
	} else if len(queryParams["name"]) > 0 {
		studios = model.Get(studios2.Studio{Studio: "%" + queryParams["name"][0] + "%"})
	} else {
		studios = model.Get(studios2.Studio{})
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

	tags := make([]tags2.Tag, 0)
	model := tags2.Initialize()

	if len(queryParams["generated_id"]) > 0 {
		genId, err := strconv.Atoi(queryParams["generated_id"][0])
		if err != nil {
			fmt.Println(err)
		}
		tags = model.Get(tags2.Tag{GeneratedID: int64(genId)})
	} else if len(queryParams["name"]) > 0 {
		tags = model.Get(tags2.Tag{Name: "%" + queryParams["name"][0] + "%"})
	} else {
		tags = model.Get(tags2.Tag{})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&tags)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func getProgress(w http.ResponseWriter, r *http.Request) {
	type progressHolder struct {
		Progress int `json:"progress"`
	}

	queryParams := r.URL.Query()

	progress := progressHolder{}
	progress.Progress = -1
	if len(queryParams["uid"]) > 0 {
		progress.Progress = manager.GetProgress(queryParams["uid"][0])
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&progress)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err := encoder.Encode(helpers.GetConfig())
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return
	case http.MethodPost:
		var config helpers.Config
		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		helpers.WriteConfig(config)

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err = encoder.Encode(helpers.GetConfig())
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	switch r.Method {
	case http.MethodPost:
		if len(queryParams["path"]) > 0 {
			err := helpers.AddPath(queryParams["path"][0])
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, "success")
			}
		}
	case http.MethodDelete:
		if len(queryParams["path"]) > 0 {
			err := helpers.RemovePath(queryParams["path"][0])
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, "success")
			}
		}
	}
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	if len(queryParams["uid"]) > 0 {
		err := manager.StopTask(queryParams["uid"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func parseMetadata(w http.ResponseWriter, r *http.Request) {
	var details struct {
		SceneId int64    `json:"generated_id,string"`
		Title   string   `json:"title"`
		Url     string   `json:"url"`
		Date    string   `json:"date"`
		Studios []string `json:"studios"`
		Actors  []string `json:"actors"`
		Tags    []string `json:"tags"`
	}
	err := json.NewDecoder(r.Body).Decode(&details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		helpers.LogError(err.Error(), component)
		return
	}

	tasks_handler.UpdateDetails(details.SceneId, details.Title, details.Date, details.Actors, details.Tags, details.Studios)
	_, err = fmt.Fprintf(w, "Success")
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func fileBrowser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var path string

	if len(queryParams["path"]) > 0 {
		path = queryParams["path"][0]
	}
	dir := GetDirectory(path)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&dir)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func queryScrapersHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	if len(queryParams["term"]) > 0 {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err := encoder.Encode(tasks_handler.GetQueryResult(queryParams["term"][0]))
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
}
