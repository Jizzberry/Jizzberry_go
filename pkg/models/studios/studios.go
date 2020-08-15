package studios

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "studios"
)

type Studio struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	Name        string `row:"studio" type:"like" json:"name"`
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

func (m Model) Create(studios []Studio) {
	// Begin Transaction
	tx, err := m.conn.Begin()
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	added := make([]string, 0)
	for _, stud := range studios {
		// Add only if value is unique
		if exists, _ := models.IsValueExists(m.conn, stud.Name, "studio", tableName); !exists {
			query, args := models.QueryBuilderCreate(stud, tableName)
			_, err := tx.Exec(query, args...)
			if err != nil {
				helpers.LogError(err.Error())
				continue
			}
			added = append(added, stud.Name)
		}
	}

	err = tx.Commit()
	if err != nil {
		helpers.LogError(err.Error())
	}

	helpers.LogInfo(fmt.Sprintf("Added studios: %v", added))
}

func (m Model) Delete(studio string) {
	query, args := models.QueryBuilderDelete(Studio{Name: studio}, tableName)
	_, err := m.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (m Model) Get(studiosQuery Studio) (allStudios []Studio) {
	query, args := models.QueryBuilderGet(studiosQuery, tableName)

	row, err := m.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return allStudios
	}

	models.GetIntoStruct(row, &allStudios)
	return
}

func (m Model) GetFromTitle(names []string) []Studio {
	fetched := make([]Studio, 0)
	for _, name := range names {
		query, args := models.QueryBuilderMatch(Studio{Name: name}, tableName)
		rows, err := m.conn.Query(query, args...)
		if err != nil {
			helpers.LogError(err.Error())
			return fetched
		}

		models.GetIntoStruct(rows, &fetched)
		fetched = removeDupl(fetched)
	}
	return fetched
}

func removeDupl(s []Studio) (list []Studio) {
	keys := make(map[int64]bool)
	for _, entry := range s {
		if _, ok := keys[entry.GeneratedID]; !ok {
			keys[entry.GeneratedID] = true
			list = append(list, entry)
		}
	}
	return
}
