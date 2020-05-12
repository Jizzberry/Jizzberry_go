package files

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
	"sync"
)

var mutexFiles = &sync.Mutex{}

type Files struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"true"`
	FileName    string `row:"file_name" type:"like"`
	FilePath    string `row:"file_path" type:"like"`
	DateCreated string `row:"date_created" type:"exact"`
	FileSize    string `row:"file_size" type:"exact"`
	Length      string `row:"length" type:"exact"`
	Tags        string `row:"tags" type:"like"`
}

type FilesModel struct {
	conn *sql.DB
}

func Initialize() *FilesModel {
	return &FilesModel{
		conn: database.GetConn(router.GetDatabase("files")),
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

	tableName := "files"
	query, args := models.QueryBuilderCreate(files, tableName)

	row, err := f.conn.Exec(query, args...)

	mutexFiles.Unlock()

	if err != nil {
		fmt.Println(err)
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
	tableName := "files"
	query, args := models.QueryBuilderDelete(files, tableName)

	if query == "" {
		return
	}

	_, err := f.conn.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func (f FilesModel) Update(files Files) {
	mutexFiles.Lock()

	if f.isEmpty() {
		mutexFiles.Unlock()
		return
	}
	tableName := "files"
	query, args := models.QueryBuilderUpdate(files, tableName)
	fmt.Println(query)

	if query == "" {
		return
	}

	_, err := f.conn.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func (f FilesModel) Get(filesQuery Files) []Files {
	tableName := "files"

	query, args := models.QueryBuilderGet(filesQuery, tableName)

	row, err := f.conn.Query(query, args...)
	if err != nil {
		fmt.Println(err)
	}

	allFiles := make([]Files, 0)
	for row.Next() {
		files := Files{}
		err := row.Scan(&files.GeneratedID, &files.FileName, &files.FilePath, &files.DateCreated, &files.FileSize, &files.Length, &files.Tags)
		if err != nil {
			fmt.Println(err)
		}
		allFiles = append(allFiles, files)
	}

	return allFiles
}

func (f FilesModel) isEmpty() bool {
	rows, err := f.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name='files'`)

	if err != nil {
		fmt.Println(err)
		return true
	}
	defer rows.Close()
	var count int

	for rows.Next() {
		rows.Scan(&count)
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
		fmt.Println(err)
	}
	var genId int64 = -1
	for fetch.Next() {
		fetch.Scan(&genId)
	}

	if genId > -1 {
		return genId, true
	}

	return -1, false
}
