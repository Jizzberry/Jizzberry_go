package files

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"strings"
	"sync"
)

var mutexFiles = &sync.Mutex{}

type Files struct {
	GeneratedID int64
	FileName    string
	FilePath    string
	DateCreated string
	FileSize    string
	Length      string
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

func (f FilesModel) Create(files *Files) int64 {
	mutexFiles.Lock()
	genId, exists := f.IsExists(files.FilePath)

	if exists {
		mutexFiles.Unlock()
		return genId
	}

	row, err := f.conn.Exec(`INSERT INTO files (file_name, file_path, date_created, file_size, length) values(?, ?, ?, ?, ?)`, files.FileName, files.FilePath, files.DateCreated, files.FileSize, files.Length)

	mutexFiles.Unlock()

	if err != nil {
		fmt.Println(err)
	}

	genID, _ := row.LastInsertId()

	defer f.Close()
	return genID
}

func (f FilesModel) Update(files Files) int64 {
	rows, err := f.conn.Exec(`UPDATE files SET file_name = ?, file_path = ? WHERE generated_id = ?`, files.FileName, files.FilePath, files.GeneratedID)
	if err != nil {
		fmt.Println(err)
	}

	genId, _ := rows.RowsAffected()
	return genId
}

func (f FilesModel) SetTags(tags []string, genId int64) {
	_, err := f.conn.Exec(`UPDATE files SET tags = ? WHERE generated_id = ?`, strings.Join(tags, ", "), genId)
	if err != nil {
		fmt.Println(err)
	}
}

func (f FilesModel) GetTags(genId int64) []string {
	rows, err := f.conn.Query(`SELECT tags FROM files WHERE generated_id = ?`, genId)
	if err != nil {
		fmt.Println(err)
	}

	tags := ""
	for rows.Next() {
		err := rows.Scan(&tags)
		if err != nil {
			fmt.Println(err)
		}
	}
	return strings.Split(tags, ", ")
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

func (f FilesModel) Get(genId int64) *Files {
	row, err := f.conn.Query(`SELECT generated_id, file_name, file_path, date_created, file_size, length FROM files WHERE generated_id = ?`, genId)
	if err != nil {
		fmt.Println(err)
	}
	files := Files{}
	for row.Next() {
		err := row.Scan(&files.GeneratedID, &files.FileName, &files.FilePath, &files.DateCreated, &files.FileSize, &files.Length)
		if err != nil {
			fmt.Println(err)
		}
	}

	return &files
}
