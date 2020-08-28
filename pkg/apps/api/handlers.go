package api

import (
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/middleware"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/studios"
	tags2 "github.com/Jizzberry/Jizzberry_go/pkg/models/tags"
	"github.com/Jizzberry/Jizzberry_go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/manager"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

type Api struct {
}

type task struct {
	Uid string `json:"uid"`
}

var decoder = schema.NewDecoder()

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
	apiRouter.HandleFunc("/organiseAll", organiseAll).Methods("POST")
	apiRouter.HandleFunc("/getMimeType", getMimeType).Methods("GET")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var file files.Files
	model := files.Initialize()
	defer model.Close()

	err := decoder.Decode(&file, r.URL.Query())

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(model.Get(file))
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func actorDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var actorDetails actor_details.ActorDetails
	model := actor_details.Initialize()
	defer model.Close()

	err := decoder.Decode(&actorDetails, r.URL.Query())

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(model.Get(actorDetails))
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func actorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var actors actor.Actor
	model := actor.Initialize()
	defer model.Close()

	err := decoder.Decode(&actors, r.URL.Query())

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(model.Get(actors))
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func scrapeActorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var actorDet actor.Actor
	err := decoder.Decode(&actorDet, r.URL.Query())

	model := actor.Initialize()
	actorDetailsModel := actor_details.Initialize()
	defer model.Close()
	defer actorDetailsModel.Close()

	tmp := model.Get(actorDet)

	var actors actor_details.ActorDetails
	if len(tmp) > 0 {
		actors = scrapers.ScrapeActor(tmp[0])
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&actors)
	if err != nil {
		helpers.LogError(err.Error())
	}

	actorDetailsModel.Create(actors)
}

func scanHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := manager.StartScan()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&task{Uid: uid})
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func scrapeListHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := make([]task, 0)
	//t = append(t, task{Uid: manager.StartScrapeActors()})
	t = append(t, task{Uid: manager.StartScrapeStudios()})
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&t)
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func studiosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var studio studios.Studio
	model := studios.Initialize()
	defer model.Close()

	err := decoder.Decode(&studio, r.URL.Query())

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(model.Get(studio))
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tag tags2.Tag
	model := tags2.Initialize()
	defer model.Close()

	err := decoder.Decode(&tag, r.URL.Query())

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(model.Get(tag))
	if err != nil {
		helpers.LogError(err.Error())
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
		helpers.LogError(err.Error())
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err := encoder.Encode(helpers.GetConfig())
		if err != nil {
			helpers.LogError(err.Error())
		}
		return
	case http.MethodPost:
		var config helpers.Config
		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = helpers.WriteConfig(config)
		if err != nil {
			helpers.LogError(err.Error())
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err = encoder.Encode(helpers.GetConfig())
		if err != nil {
			helpers.LogError(err.Error())
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
				helpers.LogError(err.Error())
				_, err := fmt.Fprintf(w, err.Error())
				if err != nil {
					helpers.LogError(err.Error())
				}
			} else {
				_, err := fmt.Fprintf(w, "success")
				if err != nil {
					helpers.LogError(err.Error())
				}
			}
		}
	case http.MethodDelete:
		if len(queryParams["path"]) > 0 {
			err := helpers.RemovePath(queryParams["path"][0])
			if err != nil {
				helpers.LogError(err.Error())
				_, err := fmt.Fprintf(w, err.Error())
				if err != nil {
					helpers.LogError(err.Error())
				}
			} else {
				_, err := fmt.Fprintf(w, "success")
				if err != nil {
					helpers.LogError(err.Error())
				}
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
	var dets tasks_handler.Details
	err := json.NewDecoder(r.Body).Decode(&dets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		helpers.LogError(err.Error())
		return
	}

	tasks_handler.UpdateDetails(dets)
	manager.StartRename(dets.SceneId)
	_, err = fmt.Fprintf(w, "Success")
	if err != nil {
		helpers.LogError(err.Error())
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
		helpers.LogError(err.Error())
	}
}

func queryScrapersHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	if len(queryParams["term"]) > 0 {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err := encoder.Encode(tasks_handler.GetQueryResult(queryParams["term"][0]))
		if err != nil {
			helpers.LogError(err.Error())
		}
	}
}

func organiseAll(w http.ResponseWriter, r *http.Request) {
	manager.OrganiseAll()
}

func getMimeType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	if len(queryParams["scene_id"]) > 0 {
		sceneId, err := strconv.ParseInt(queryParams["scene_id"][0], 10, 64)
		if err != nil {
			helpers.LogError(err.Error())
		}
		model := files.Initialize()
		defer model.Close()

		file := model.Get(files.Files{SceneID: sceneId})
		if len(file) > 0 {
			_, err = fmt.Fprintf(w, file[0].ExtraCodec)
			if err != nil {
				helpers.LogError(err.Error())
			}
		}
	}
}
