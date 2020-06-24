package router

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"io/ioutil"
	"path/filepath"
)

func GetDatabase(table string) string {
	var databaseDir = helpers.GetWorkingDirectory() + "/assets/database"
	switch table {
	case "files":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "actor_details":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "tags":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "studios":
		return filepath.FromSlash(databaseDir + "/studios.db")

	case "actors":
		return filepath.FromSlash(databaseDir + "/actors.db")

	case "auth":
		return filepath.FromSlash(databaseDir + "/auth.db")
	}

	file, _ := ioutil.TempFile(databaseDir, "/tmp.db")
	return file.Name()
}

func GetJson(name string) string {
	var databaseDir = helpers.GetWorkingDirectory() + "/assets/json"
	switch name {
	case "actorsRelation":
		return filepath.FromSlash(databaseDir + "/actorsRelation.json")

	case "studiosRelation":
		return filepath.FromSlash(databaseDir + "/studiosRelation.json")

	case "tagsRelation":
		return filepath.FromSlash(databaseDir + "/tagsRelation.json")
	}
	file, _ := ioutil.TempFile(databaseDir, "/tmp.json")
	return file.Name()
}
