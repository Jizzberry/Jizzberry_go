package actor

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
)

type Actor struct {
	GeneratedID int64
	Name        string
	UrlID       string
	Website     string
}

type Actors struct {
	Actors []Actor
}

type ActorsModel struct {
	conn *sql.DB
}

func Initialize() *ActorsModel {
	return &ActorsModel{
		conn: database.GetConn(router.GetDatabase("actors")),
	}
}

func (a ActorsModel) Close() {
	a.conn.Close()
}

func (a ActorsModel) Create(actors Actors) {
	tx, err := a.conn.Begin()

	if err != nil {
		fmt.Println(err)
		return
	}

	if a.isEmpty() {
		database.RunMigrations()
	}

	for _, act := range actors.Actors {
		_, err := tx.Exec(`INSERT INTO actors (name, website, urlid) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT 1 FROM actors WHERE name = ?)`, act.Name, act.Website, act.UrlID, act.Name)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
	}

	defer a.Close()
}

func (a ActorsModel) GetExact(name string) Actor {

	rows, err := a.conn.Query(`SELECT generated_id, name, website, urlid FROM actors WHERE name = ?`, name)
	if err != nil {
		fmt.Println(err)
	}

	actor := Actor{Name: name}
	for rows.Next() {
		rows.Scan(&actor.GeneratedID, &actor.Name, &actor.Website, &actor.UrlID)
	}
	return actor
}

func (a ActorsModel) Get(names []string) [][]Actor {
	final := make([][]Actor, len(names))
	for i, name := range names {
		rows, err := a.conn.Query(`SELECT generated_id, name, website, urlid FROM actors WHERE (name LIKE ? COLLATE NOCASE) 
                                                         OR (replace(name, ' ', '') LIKE ? COLLATE NOCASE)`, "%"+name+"%", name)
		if err != nil {
			fmt.Println(err)
		}

		fetched := make([]Actor, 0)

		for rows.Next() {
			var actor = Actor{}
			err := rows.Scan(&actor.GeneratedID, &actor.Name, &actor.Website, &actor.UrlID)
			if err != nil {
				fmt.Println(err)
			}
			fetched = append(fetched, actor)
		}
		final[i] = fetched
	}
	return final
}

func (a ActorsModel) isEmpty() bool {
	rows, err := a.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name='actor'`)

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
