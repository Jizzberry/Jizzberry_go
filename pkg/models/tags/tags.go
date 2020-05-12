package tags

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
)

type Tags struct {
	GeneratedID int64
	Tags        string
}

type TagsModel struct {
	conn *sql.DB
}

func Initialize() *TagsModel {
	return &TagsModel{
		conn: database.GetConn(router.GetDatabase("tags")),
	}
}

func (t TagsModel) Close() {
	t.conn.Close()
}

func (t TagsModel) isEmpty() bool {
	rows, err := t.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name='tags'`)

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

func (t TagsModel) IsExists(filePath string) (int64, bool) {
	if t.isEmpty() {
		database.RunMigrations()
		return -1, false
	}

	fetch, err := t.conn.Query(`SELECT generated_id FROM tags WHERE tag = ?`, filePath)
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

func (t TagsModel) Create(tags *Tags) int64 {
	genId, exists := t.IsExists(tags.Tags)

	if exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(tags, "tags")

	row, err := t.conn.Exec(query, args...)

	if err != nil {
		fmt.Println(err)
	}

	genID, _ := row.LastInsertId()

	defer t.Close()
	return genID
}

func (t TagsModel) Delete(tag string) {
	_, err := t.conn.Exec(`DELETE FROM tags WHERE tag = ?`, tag)

	if err != nil {
		fmt.Println(err)
	}
}
