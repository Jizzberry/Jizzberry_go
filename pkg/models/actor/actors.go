package actor

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
)

const (
	tableName = "actors"
	component = "actorsModel"
)

type Actor struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"true" json:"generated_id"`
	Name        string `row:"name" type:"like" json:"name"`
	UrlID       string `row:"urlid" type:"exact" json:"urlid"`
	Website     string `row:"website" type:"exact" json:"website"`
}

type ActorsModel struct {
	conn *sql.DB
}

func Initialize() *ActorsModel {
	return &ActorsModel{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (a ActorsModel) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (a ActorsModel) Create(actors []Actor) {
	tx, err := a.conn.Begin()

	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	if a.isEmpty() {
		database.RunMigrations()
	}

	for _, act := range actors {
		_, err := tx.Exec(`INSERT INTO actors (name, website, urlid) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT 1 FROM actors WHERE name = ?)`, act.Name, act.Website, act.UrlID, act.Name)
		if err != nil {
			helpers.LogError(err.Error(), component)
			tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		helpers.LogError(err.Error(), component)
		tx.Rollback()
	}

	defer a.Close()
}

func (a ActorsModel) GetExact(name string) Actor {
	actor := Actor{Name: name}

	rows, err := a.conn.Query(`SELECT generated_id, name, website, urlid FROM actors WHERE name = ?`, name)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return actor
	}

	for rows.Next() {
		err := rows.Scan(&actor.GeneratedID, &actor.Name, &actor.Website, &actor.UrlID)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
	return actor
}

func (a ActorsModel) GetFromTitle(names []string) []Actor {
	fetched := make([]Actor, 0)
	for _, name := range names {
		rows, err := a.conn.Query(`SELECT generated_id, name, website, urlid FROM actors WHERE (name LIKE ? COLLATE NOCASE) 
                                                         OR (replace(name, ' ', '') LIKE ? COLLATE NOCASE)`, "%"+name+"%", name)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return fetched
		}

		for rows.Next() {
			var actor = Actor{}
			err := rows.Scan(&actor.GeneratedID, &actor.Name, &actor.Website, &actor.UrlID)
			if err != nil {
				helpers.LogError(err.Error(), component)
			}

			if !containsActors(fetched, actor) {
				fetched = append(fetched, actor)
			}
		}
	}
	return fetched
}

func (a ActorsModel) Get(actor Actor) []Actor {
	query, args := models.QueryBuilderGet(actor, tableName)

	rows, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	allActors := make([]Actor, 0)
	for rows.Next() {
		actor := Actor{}
		err := rows.Scan(&actor.GeneratedID, &actor.Name, &actor.UrlID, &actor.Website)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allActors = append(allActors, actor)
	}

	return allActors
}

func (a ActorsModel) isEmpty() bool {
	rows, err := a.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

	if err != nil {
		helpers.LogError(err.Error(), component)
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

func containsActors(s []Actor, e Actor) bool {
	for _, a := range s {
		if a.GeneratedID == e.GeneratedID {
			return true
		}
	}
	return false
}
