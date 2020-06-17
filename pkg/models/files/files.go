package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/studios"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/tags"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

var mutexFiles = &sync.Mutex{}

const (
	tableName = "files"
	component = "filesModel"
)

type Files struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id,string"`
	FileName    string `row:"file_name" type:"like" json:"file_name"`
	FilePath    string `row:"file_path" type:"like" json:"file_path"`
	DateCreated string `row:"date_created" type:"exact" json:"date_created"`
	FileSize    string `row:"file_size" type:"exact" json:"file_size"`
	Length      string `row:"length" type:"exact" json:"length"`
	Tags        string `row:"tags" type:"like" json:"tags"`
	Studios     string `row:"studios" type:"like" json:"studios"`
	Actors      string `row:"actors" type:"like" json:"actors"`
	URL         string `row:"url" type:"exact" json:"url"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (f Model) Close() {
	err := f.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (f Model) Create(files Files) int64 {
	mutexFiles.Lock()
	genId, exists := f.IsExists(files.FilePath)

	if exists {
		mutexFiles.Unlock()
		return genId
	}

	query, args := models.QueryBuilderCreate(files, tableName)

	row, err := f.conn.Exec(query, args...)

	mutexFiles.Unlock()

	if err != nil {
		helpers.LogError(err.Error(), component)
		return 0
	}

	genID, err := row.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	setActorRelation(genID, files.Actors)
	setStudioRelation(genID, files.Studios)
	setTagRelation(genID, files.Tags)
	return genID
}

func (f Model) Delete(files Files) {
	mutexFiles.Lock()

	if f.isEmpty() {
		mutexFiles.Unlock()
		return
	}

	query, args := models.QueryBuilderDelete(files, tableName)

	if query == "" {
		return
	}

	_, err := f.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		mutexFiles.Unlock()
		return
	}
	setActorRelation(files.GeneratedID, "")
	setStudioRelation(files.GeneratedID, "")
	setTagRelation(files.GeneratedID, "")
	mutexFiles.Unlock()
}

func (f Model) Update(files Files) {
	mutexFiles.Lock()

	if f.isEmpty() {
		mutexFiles.Unlock()
		return
	}

	query, args := models.QueryBuilderUpdate(files, tableName)

	if query == "" {
		return
	}

	_, err := f.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		mutexFiles.Unlock()
		return
	}

	setActorRelation(files.GeneratedID, files.Actors)
	setStudioRelation(files.GeneratedID, files.Studios)
	setTagRelation(files.GeneratedID, files.Tags)
	mutexFiles.Unlock()

}

func (f Model) Get(filesQuery Files) []Files {
	query, args := models.QueryBuilderGet(filesQuery, tableName)
	allFiles := make([]Files, 0)

	row, err := f.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return allFiles
	}

	for row.Next() {
		files := Files{}
		err := row.Scan(&files.GeneratedID, &files.FileName, &files.FilePath, &files.DateCreated, &files.FileSize, &files.Length, &files.Tags, &files.Studios, &files.Actors, &files.URL)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allFiles = append(allFiles, files)
	}

	return allFiles
}

func (f Model) isEmpty() bool {
	rows, err := f.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

	if err != nil {
		helpers.LogError(err.Error(), component)
		return true
	}
	var count int

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if count < 0 {
		return true
	}
	return false
}

func (f Model) IsExists(filePath string) (int64, bool) {
	if f.isEmpty() {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return -1, false
	}

	fetch, err := f.conn.Query(`SELECT generated_id FROM files WHERE file_path = ?`, filePath)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return -1, false
	}
	var genId int64 = -1
	for fetch.Next() {
		err := fetch.Scan(&genId)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if genId > -1 {
		return genId, true
	}

	return -1, false
}

func setActorRelation(genId int64, actors string) {

	actorsSli := strings.Split(actors, ", ")

	jsonFile := readJson(router.GetJson("actorsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	deleteRelation(&relation, strconv.FormatInt(genId, 10))
	if relation != nil {
		if actors != "" {
			actorsModel := actor.Initialize()
			defer actorsModel.Close()

			for _, a := range actorsSli {
				tmp := actorsModel.GetExact(a)
				relation[strconv.FormatInt(tmp.GeneratedID, 10)] = append(relation[strconv.FormatInt(tmp.GeneratedID, 10)], strconv.FormatInt(genId, 10))
			}
		}
		writeJson(jsonFile, relation)
	}
}

func GetActorRelations(ActorID string) []string {
	jsonFile := readJson(router.GetJson("actorsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	if val, ok := relation[ActorID]; ok {
		return val
	}
	return nil
}

func setStudioRelation(genId int64, studio string) {
	split := strings.Split(studio, ", ")

	jsonFile := readJson(router.GetJson("studiosRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	deleteRelation(&relation, strconv.FormatInt(genId, 10))

	if relation != nil {
		if studio != "" {
			studiosModel := studios.Initialize()
			defer studiosModel.Close()
			for _, s := range split {
				tmp := studiosModel.Get(studios.Studio{Studio: s})[0]
				relation[strconv.FormatInt(tmp.GeneratedID, 10)] = append(relation[strconv.FormatInt(tmp.GeneratedID, 10)], strconv.FormatInt(genId, 10))
			}
		}

		writeJson(jsonFile, relation)
	}

}

func setTagRelation(genId int64, tag string) {
	split := strings.Split(tag, ", ")

	jsonFile := readJson(router.GetJson("tagsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	deleteRelation(&relation, strconv.FormatInt(genId, 10))

	if relation != nil {
		if tag != "" {
			tagsModel := tags.Initialize()
			defer tagsModel.Close()

			for _, t := range split {
				tmp := tagsModel.Get(tags.Tag{Name: t})[0]
				relation[strconv.FormatInt(tmp.GeneratedID, 10)] = append(relation[strconv.FormatInt(tmp.GeneratedID, 10)], strconv.FormatInt(genId, 10))
			}
		}

		writeJson(jsonFile, relation)
	}

}

func deleteRelation(relations *map[string][]string, genID string) {
	if relations != nil {
		for key, value := range *relations {
			for _, v := range value {
				if v == genID {
					delete(*relations, key)
				}
			}
		}
	}
}

func readJson(filename string) *os.File {
	jsonFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return nil
	}
	return jsonFile
}

func parseJson(file *os.File) map[string][]string {

	byteValue, _ := ioutil.ReadAll(file)
	relation := make(map[string][]string)

	if len(byteValue) > 0 {
		err := json.Unmarshal(byteValue, &relation)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
	return relation
}

func writeJson(file *os.File, relation map[string][]string) {
	bytes, err := json.Marshal(relation)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	_, err = file.WriteAt(bytes, 0)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}
