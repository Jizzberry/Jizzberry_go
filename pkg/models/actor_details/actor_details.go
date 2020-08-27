package actor_details

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "actor_details"
)

type ActorDetails struct {
	GeneratedId   int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	ActorId       int64  `row:"actor_id" type:"exact" json:"actor_id"`
	Name          string `row:"name" type:"like" json:"name"`
	Birthday      string `row:"born" type:"like" json:"birthday"`
	Birthplace    string `row:"birthplace" type:"like" json:"birthplace"`
	Height        string `row:"height" type:"like" json:"height"`
	Weight        string `row:"weight" type:"like" json:"weight"`
	ThumbnailPath string `row:"thumbnail" type:"exact" json:"thumbnail"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: models.GetConn(tableName),
	}
}

func (a Model) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (a Model) Create(details ActorDetails) int64 {
	if ok, gen := models.IsValueExists(a.conn, details.ActorId, "actor_id", tableName); ok {
		return gen
	}

	query, args := models.QueryBuilderCreate(details, tableName)

	row, err := a.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error())
		return 0
	}

	genID, _ := row.LastInsertId()
	return genID
}

func (a Model) Delete(details ActorDetails) {
	query, args := models.QueryBuilderDelete(details, tableName)

	_, err := a.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (a Model) Get(d ActorDetails) (allDetails []ActorDetails) {
	query, args := models.QueryBuilderGet(d, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &allDetails)
	return
}

func (a Model) IsExists(actorId int64) bool {
	exists, _ := models.IsValueExists(a.conn, actorId, "actor_id", tableName)
	return exists
}
