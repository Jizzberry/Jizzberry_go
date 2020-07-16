package actor

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "actors"
	component = "actorsModel"
)

type Actor struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	Name        string `row:"name" type:"like" json:"name"`
	UrlID       string `row:"urlid" type:"like" json:"urlid"`
	Website     string `row:"website" type:"like" json:"website"`
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
		helpers.LogError(err.Error(), component)
	}
}

func (a Model) Create(actors []Actor) {
	// Start transaction
	tx, err := a.conn.Begin()

	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	// Commit all actors from transaction
	added := make([]string, 0)
	for _, act := range actors {

		// Add only if value is unique
		if exist, _ := models.IsValueExists(a.conn, act.Name, "name", tableName); !exist {
			query, args := models.QueryBuilderCreate(act, tableName)
			_, err := tx.Exec(query, args...)
			if err != nil {
				helpers.LogError(err.Error(), component)
				continue
			}
			added = append(added, act.Name)
		}
	}

	err = tx.Commit()
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	helpers.LogInfo(fmt.Sprintf("Added actors: %v", added), component)
}

// Match actors from array of words
func (a Model) GetFromTitle(names []string) (fetched []Actor) {
	for _, name := range names {
		query, args := models.QueryBuilderMatch(Actor{Name: name}, tableName)
		rows, err := a.conn.Query(query, args...)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return fetched
		}

		models.GetIntoStruct(rows, &fetched)
		fetched = removeDupl(fetched)
	}
	return
}

func (a Model) Get(actor Actor) (allActors []Actor) {
	query, args := models.QueryBuilderGet(actor, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	models.GetIntoStruct(row, &allActors)
	return
}

func removeDupl(s []Actor) (list []Actor) {
	keys := make(map[int64]bool)
	for _, entry := range s {
		if _, ok := keys[entry.GeneratedID]; !ok {
			keys[entry.GeneratedID] = true
			list = append(list, entry)
		}
	}
	return
}
