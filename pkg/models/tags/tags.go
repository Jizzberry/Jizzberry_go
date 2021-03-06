package tags

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "tags"
)

type Tag struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id" schema:"generated_id"`
	Name        string `row:"tag" type:"like" json:"name" schema:"name"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: models.GetConn(tableName),
	}
}

func (t Model) Close() {
	err := t.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (t Model) Create(tags Tag) int64 {
	exists, genId := models.IsValueExists(t.conn, tags.Name, "tag", tableName)

	if exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(tags, tableName)
	row, err := t.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error())
		return 0
	}

	genID, _ := row.LastInsertId()

	return genID
}

func (t Model) Delete(tag string) {
	query, args := models.QueryBuilderDelete(Tag{Name: tag}, tableName)
	_, err := t.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (t Model) Get(tagsQuery Tag) (allTags []Tag) {
	query, args := models.QueryBuilderGet(tagsQuery, tableName)

	row, err := t.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return allTags
	}

	models.GetIntoStruct(row, &allTags)
	return
}
