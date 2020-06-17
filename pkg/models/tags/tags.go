package tags

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
)

const (
	tableName = "tags"
	component = "tagsModel"
)

type Tag struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	Name        string `row:"tag" type:"like" json:"generated_id"`
	Count       int64  `row:"count" type:"exact" json:"count"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (t Model) Close() {
	err := t.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (t Model) isEmpty() bool {
	rows, err := t.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

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

func (t Model) IsExists(filePath string) (int64, bool) {
	if t.isEmpty() {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return -1, false
	}

	fetch, err := t.conn.Query(`SELECT generated_id FROM tags WHERE tag = ?`, filePath)
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

func (t Model) Create(tags Tag) int64 {
	genId, exists := t.IsExists(tags.Name)

	if exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(tags, tableName)

	row, err := t.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error(), component)
		return 0
	}

	genID, _ := row.LastInsertId()

	return genID
}

func (t Model) Delete(tag string) {
	_, err := t.conn.Exec(`DELETE FROM tags WHERE tag = ?`, tag)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (t Model) Get(tagsQuery Tag) []Tag {
	query, args := models.QueryBuilderGet(tagsQuery, tableName)
	allTags := make([]Tag, 0)

	row, err := t.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return allTags
	}

	for row.Next() {
		tag := Tag{}
		err := row.Scan(&tag.GeneratedID, &tag.Name, &tag.Count)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allTags = append(allTags, tag)
	}

	return allTags
}
