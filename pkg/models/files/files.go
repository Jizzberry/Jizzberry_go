package files

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "files"
	component = "filesModel"
)

type Files struct {
	GeneratedID   int64   `row:"generated_id" type:"exact" pk:"auto" json:"generated_id,string"`
	FileName      string  `row:"file_name" type:"like" json:"file_name"`
	FilePath      string  `row:"file_path" type:"like" json:"file_path"`
	DateCreated   string  `row:"date_created" type:"exact" json:"date_created"`
	FileSize      string  `row:"file_size" type:"exact" json:"file_size"`
	ThumbnailPath string  `row:"thumbnail" type:"exact" json:"thumbnail"`
	Symlinks      string  `row:"symlinks" type:"exact" json:"symlinks"`
	Tags          string  `row:"tags" type:"like" json:"tags"`
	Studios       string  `row:"studios" type:"like" json:"studios"`
	Actors        string  `row:"actors" type:"like" json:"actors"`
	URL           string  `row:"url" type:"exact" json:"url"`
	Length        float64 `row:"length" type:"exact" json:"length"`
	Format        string  `row:"format" type:"like" json:"format"`
	Video0Codec   string  `row:"video0codec" type:"like" json:"video0_codec"`
	Audio0Codec   string  `row:"audio0codec" type:"like" json:"audio0_codec"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: models.GetConn(tableName),
	}
}

func (m Model) Close() {
	err := m.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (m Model) Create(files Files) int64 {
	if exists, genId := models.IsValueExists(m.conn, files.FilePath, "file_path", tableName); exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(files, tableName)
	row, err := m.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return 0
	}

	genID, err := row.LastInsertId()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	setAllRelations(genID, files.Actors, files.Studios, files.Tags)
	return genID
}

func (m Model) Delete(files Files) {

	query, args := models.QueryBuilderDelete(files, tableName)

	if query == "" {
		return
	}

	_, err := m.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}
	setAllRelations(files.GeneratedID, "", "", "")
}

func (m Model) Update(files Files) {
	query, args := models.QueryBuilderUpdate(files, tableName)

	if query == "" {
		return
	}

	_, err := m.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}
	setAllRelations(files.GeneratedID, files.Actors, files.Studios, files.Tags)
}

func (m Model) Get(filesQuery Files) (allFiles []Files) {
	query, args := models.QueryBuilderGet(filesQuery, tableName)

	row, err := m.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	models.GetIntoStruct(row, &allFiles)
	return
}

func (m Model) IsExists(filepath string) bool {
	exists, _ := models.IsValueExists(m.conn, filepath, "file_path", tableName)
	return exists
}
