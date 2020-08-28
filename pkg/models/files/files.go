package files

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "files"
)

type Files struct {
	SceneID       int64   `row:"scene_id" type:"exact" pk:"auto" json:"scene_id" schema:"scene_id"`
	FileName      string  `row:"file_name" type:"like" json:"file_name" schema:"file_name"`
	FilePath      string  `row:"file_path" type:"like" json:"file_path" schema:"file_path"`
	DateCreated   string  `row:"date_created" type:"exact" json:"date_created" schema:"date_created"`
	FileSize      string  `row:"file_size" type:"exact" json:"file_size" schema:"file_size"`
	ThumbnailPath string  `row:"thumbnail" type:"exact" json:"thumbnail" schema:"thumbnail"`
	Symlinks      string  `row:"symlinks" type:"exact" json:"symlinks" schema:"symlinks"`
	Tags          string  `row:"tags" type:"like" json:"tags" schema:"tags"`
	Studios       string  `row:"studios" type:"like" json:"studios" schema:"studios"`
	Actors        string  `row:"actors" type:"like" json:"actors" schema:"actors"`
	URL           string  `row:"url" type:"exact" json:"url" schema:"url"`
	Length        float64 `row:"length" type:"exact" json:"length" schema:"length"`
	Format        string  `row:"format" type:"like" json:"format" schema:"format"`
	Video0Codec   string  `row:"video0codec" type:"like" json:"video0_codec" schema:"video0_codec"`
	Audio0Codec   string  `row:"audio0codec" type:"like" json:"audio0_codec" schema:"audio0_codec"`
	ExtraCodec    string  `row:"extra_codec" type:"exact" json:"extra_codec" schema:"extra_codec"`
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
		helpers.LogError(err.Error())
	}
}

func (m Model) Create(files Files) int64 {
	if exists, genId := models.IsValueExists(m.conn, files.FilePath, "file_path", tableName); exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(files, tableName)
	row, err := m.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return 0
	}

	genID, err := row.LastInsertId()
	if err != nil {
		helpers.LogError(err.Error())
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
		helpers.LogError(err.Error())
		return
	}
	setAllRelations(files.SceneID, "", "", "")
}

func (m Model) Update(files Files) {
	query, args := models.QueryBuilderUpdate(files, tableName)

	if query == "" {
		return
	}

	_, err := m.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	setAllRelations(files.SceneID, files.Actors, files.Studios, files.Tags)
}

func (m Model) Get(filesQuery Files) (allFiles []Files) {
	query, args := models.QueryBuilderGet(filesQuery, tableName)

	row, err := m.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &allFiles)
	return
}

func (m Model) IsExists(filepath string) bool {
	exists, _ := models.IsValueExists(m.conn, filepath, "file_path", tableName)
	return exists
}
