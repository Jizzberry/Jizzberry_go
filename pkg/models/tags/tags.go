package tags

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/logging"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
)

const (
	tableName = "tags"
	component = "tagsModel"
)

type Tags struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"true" json:"generated_id"`
	Tags        string `row:"tags" type:"like" json:"generated_id"`
}

type TagsModel struct {
	conn *sql.DB
}

func Initialize() *TagsModel {
	return &TagsModel{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (t TagsModel) Close() {
	t.conn.Close()
}

func (t TagsModel) isEmpty() bool {
	rows, err := t.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

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

func (t TagsModel) IsExists(filePath string) (int64, bool) {
	if t.isEmpty() {
		database.RunMigrations()
		return -1, false
	}

	fetch, err := t.conn.Query(`SELECT generated_id FROM tags WHERE tag = ?`, filePath)
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

func (t TagsModel) Create(tags *Tags) int64 {
	genId, exists := t.IsExists(tags.Tags)

	if exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(tags, tableName)

	row, err := t.conn.Exec(query, args...)

	if err != nil {
		logging.LogError(err.Error(), component)
		return 0
	}

	genID, _ := row.LastInsertId()

	defer t.Close()
	return genID
}

func (t TagsModel) Delete(tag string) {
	_, err := t.conn.Exec(`DELETE FROM tags WHERE tag = ?`, tag)

	if err != nil {
		logging.LogError(err.Error(), component)
	}
}
