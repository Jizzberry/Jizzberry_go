package actor_details

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"sync"
)

var mutexDetails = &sync.Mutex{}

type ActorDetails struct {
	SceneId    int64
	ActorId    int64
	Name       string
	Birthday   string
	Birthplace string
	Height     string
	Weight     string
}

type ActorDetailsModel struct {
	conn *sql.DB
}

func Initialize() *ActorDetailsModel {
	return &ActorDetailsModel{
		conn: database.GetConn(router.GetDatabase("actor_details")),
	}
}

func (a ActorDetailsModel) Close() {
	a.conn.Close()
}

func (a ActorDetailsModel) Create(details *ActorDetails) int64 {
	mutexDetails.Lock()

	if a.isEmpty() {
		database.RunMigrations()
	}

	row, err := a.conn.Exec(`INSERT INTO actor_details (actor_id, scene_id, name, born, birthplace, height, weight) values (?, ?, ?, ?, ?, ?, ?)`,
		details.ActorId, details.SceneId, details.Name, details.Birthday, details.Birthplace, details.Height, details.Weight)

	mutexDetails.Unlock()

	if err != nil {
		fmt.Println(err)
	}

	genID, _ := row.LastInsertId()
	defer a.Close()
	return genID
}

func (a ActorDetailsModel) Delete(sceneId int64) {
	if a.isEmpty() {
		database.RunMigrations()
		return
	}

	_, err := a.conn.Exec(`DELETE FROM actor_details WHERE scene_id = ?`, sceneId)

	if err != nil {
		fmt.Println(err)
	}
}

func (a ActorDetailsModel) Get(actorId int64) *ActorDetails {
	mutexDetails.Lock()
	row, err := a.conn.Query(`SELECT actor_id, name, born, birthplace, height, weight FROM actor_details WHERE actor_id = ?`, actorId)
	if err != nil {
		fmt.Println(err)
	}

	details := ActorDetails{}
	for row.Next() {
		err := row.Scan(&details.ActorId, &details.Name, &details.Birthday, &details.Birthplace, &details.Height, &details.Weight)
		if err != nil {
			fmt.Println(err)
		}
	}
	return &details
}

func (a ActorDetailsModel) IsExists(actorId int64) bool {
	rows, err := a.conn.Query(`SELECT actor_id FROM actor_details WHERE actor_id=?`, actorId)

	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows.Close()
	if rows.NextResultSet() {
		return true
	}
	return false
}

func (a ActorDetailsModel) isEmpty() bool {
	rows, err := a.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name='actor_details'`)

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
