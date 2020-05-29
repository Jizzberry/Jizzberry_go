package router

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"io/ioutil"
	"path/filepath"
)

func GetDatabase(table string) string {
	databaseDir := helpers.GetWorkingDirectory() + "/assets/database"
	switch table {
	case "files":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "actor_details":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "tags":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "studios":
		return filepath.FromSlash(databaseDir + "/jizzberry_data.db")

	case "actors":
		return filepath.FromSlash(databaseDir + "/actors.db")

	case "auth":
		return filepath.FromSlash(databaseDir + "/auth.db")
	}

	file, _ := ioutil.TempFile(databaseDir, "/tmp.db")
	return file.Name()
}
