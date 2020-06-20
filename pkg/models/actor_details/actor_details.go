package actor_details

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/database"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
	"sync"
)

const (
	tableName = "actor_details"
	component = "actorDetailsModel"
)

var mutexDetails = &sync.Mutex{}

type ActorDetails struct {
	GeneratedId int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	ActorId     int64  `row:"actor_id" type:"exact" json:"actor_id"`
	Name        string `row:"name" type:"like" json:"name"`
	Birthday    string `row:"born" type:"like" json:"birthday"`
	Birthplace  string `row:"birthplace" type:"like" json:"birthplace"`
	Height      string `row:"height" type:"exact" json:"height"`
	Weight      string `row:"weight" type:"exact" json:"weight"`
}

type ActorDetailsModel struct {
	conn *sql.DB
}

func Initialize() *ActorDetailsModel {
	return &ActorDetailsModel{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (a ActorDetailsModel) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (a ActorDetailsModel) Create(details ActorDetails) int64 {
	mutexDetails.Lock()

	if a.isEmpty() {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if ok, gen := a.IsExists(details.ActorId); ok {
		mutexDetails.Unlock()
		return gen
	}

	query, args := models.QueryBuilderCreate(details, tableName)

	row, err := a.conn.Exec(query, args...)

	mutexDetails.Unlock()

	if err != nil {
		mutexDetails.Unlock()
		helpers.LogError(err.Error(), component)
		return 0
	}

	genID, _ := row.LastInsertId()
	return genID
}

func (a ActorDetailsModel) Delete(details ActorDetails) {
	if a.isEmpty() {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return
	}

	query, args := models.QueryBuilderDelete(details, tableName)

	_, err := a.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (a ActorDetailsModel) Get(d ActorDetails) []ActorDetails {
	mutexDetails.Lock()

	allDetails := make([]ActorDetails, 0)

	query, args := models.QueryBuilderGet(d, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return allDetails
	}

	for row.Next() {
		details := ActorDetails{}
		err := row.Scan(&details.GeneratedId, &details.ActorId, &details.Name, &details.Birthday, &details.Birthplace, &details.Height, &details.Weight)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allDetails = append(allDetails, details)
	}
	mutexDetails.Unlock()
	return allDetails
}

func (a ActorDetailsModel) IsExists(actorId int64) (bool, int64) {
	rows, err := a.conn.Query(`SELECT actor_id FROM actor_details WHERE actor_id=?`, actorId)

	if err != nil {
		helpers.LogError(err.Error(), component)
		return false, -1
	}

	var genId int64 = -1
	for rows.Next() {
		err := rows.Scan(&genId)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if genId > -1 {
		return true, genId
	}

	return false, -1
}

func (a ActorDetailsModel) isEmpty() bool {
	rows, err := a.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name='actor_details'`)

	if err != nil {
		fmt.Println(err)
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
