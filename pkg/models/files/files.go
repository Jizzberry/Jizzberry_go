package files

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/logging"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
	"sync"
)

var mutexFiles = &sync.Mutex{}

const (
	tableName = "files"
	component = "filesModel"
)

type Files struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"true" json:"generated_id"`
	FileName    string `row:"file_name" type:"like" json:"file_name"`
	FilePath    string `row:"file_path" type:"like" json:"file_path"`
	DateCreated string `row:"date_created" type:"exact" json:"date_created"`
	FileSize    string `row:"file_size" type:"exact" json:"file_size"`
	Length      string `row:"length" type:"exact" json:"length"`
	Tags        string `row:"tags" type:"like" json:"tags"`
}

type FilesModel struct {
	conn *sql.DB
}

func Initialize() *FilesModel {
	return &FilesModel{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (f FilesModel) Close() {
	f.conn.Close()
}

func (f FilesModel) Create(files Files) int64 {
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
		logging.LogError(err.Error(), component)
		return 0
	}

	genID, _ := row.LastInsertId()

	defer f.Close()
	return genID
}

func (f FilesModel) Delete(files Files) {
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
		logging.LogError(err.Error(), component)
	}
}

func (f FilesModel) Update(files Files) {
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
		logging.LogError(err.Error(), component)
	}
}

func (f FilesModel) Get(filesQuery Files) []Files {
	query, args := models.QueryBuilderGet(filesQuery, tableName)
	allFiles := make([]Files, 0)

	row, err := f.conn.Query(query, args...)
	if err != nil {
		logging.LogError(err.Error(), component)
		return allFiles
	}

	for row.Next() {
		files := Files{}
		err := row.Scan(&files.GeneratedID, &files.FileName, &files.FilePath, &files.DateCreated, &files.FileSize, &files.Length, &files.Tags)
		if err != nil {
			logging.LogError(err.Error(), component)
		}
		allFiles = append(allFiles, files)
	}

	return allFiles
}

func (f FilesModel) isEmpty() bool {
	rows, err := f.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

	if err != nil {
		logging.LogError(err.Error(), component)
		return true
	}
	var count int

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			logging.LogError(err.Error(), component)
		}
	}

	if count < 0 {
		return true
	}
	return false
}

func (f FilesModel) IsExists(filePath string) (int64, bool) {
	if f.isEmpty() {
		database.RunMigrations()
		return -1, false
	}

	fetch, err := f.conn.Query(`SELECT generated_id FROM files WHERE file_path = ?`, filePath)
	if err != nil {
		logging.LogError(err.Error(), component)
		return -1, false
	}
	var genId int64 = -1
	for fetch.Next() {
		err := fetch.Scan(&genId)
		if err != nil {
			logging.LogError(err.Error(), component)
		}
	}

	if genId > -1 {
		return genId, true
	}

	return -1, false
}
