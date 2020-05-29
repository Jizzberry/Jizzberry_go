package studios

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
)

const (
	tableName = "studios"
	component = "studiosModel"
)

type Studios struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"true" json:"generated_id"`
	Studio      string `row:"studio" type:"like" json:"generated_id"`
}

type StudiosModel struct {
	conn *sql.DB
}

func Initialize() *StudiosModel {
	return &StudiosModel{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (s StudiosModel) Close() {
	s.conn.Close()
}

func (s StudiosModel) isEmpty() bool {
	rows, err := s.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

	if err != nil {
		helpers.LogError(err.Error(), component)
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

func (s StudiosModel) IsExists(studio string) (int64, bool) {
	if s.isEmpty() {
		database.RunMigrations()
		return -1, false
	}

	fetch, err := s.conn.Query(`SELECT generated_id FROM studios WHERE studio = ?`, studio)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return -1, false
	}
	var genId int64 = -1
	for fetch.Next() {
		err := fetch.Scan(&genId)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if genId > -1 {
		return genId, true
	}

	return -1, false
}

func (s StudiosModel) Create(studios Studios) int64 {
	genId, exists := s.IsExists(studios.Studio)

	if exists {
		return genId
	}

	query, args := models.QueryBuilderCreate(studios, tableName)

	row, err := s.conn.Exec(query, args...)

	if err != nil {
		helpers.LogError(err.Error(), component)
		return 0
	}

	genID, _ := row.LastInsertId()

	defer s.Close()
	return genID
}

func (s StudiosModel) Delete(studio string) {
	_, err := s.conn.Exec(`DELETE FROM tags WHERE tag = ?`, studio)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (s StudiosModel) Get(studiosQuery Studios) []Studios {
	query, args := models.QueryBuilderGet(studiosQuery, tableName)
	allStudios := make([]Studios, 0)

	row, err := s.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return allStudios
	}

	for row.Next() {
		studios := Studios{}
		err := row.Scan(&studios.GeneratedID, &studios.Studio)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allStudios = append(allStudios, studios)
	}

	return allStudios
}
