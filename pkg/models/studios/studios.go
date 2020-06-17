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

type Studio struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	Studio      string `row:"studio" type:"like" json:"generated_id"`
	Count       int64  `row:"count" type:"exact" json:"count"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (s Model) Close() {
	err := s.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (s Model) isEmpty() bool {
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

func (s Model) IsExists(studio string) (int64, bool) {
	if s.isEmpty() {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
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

func (s Model) Create(studios []Studio) {
	tx, err := s.conn.Begin()

	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	for _, stud := range studios {
		_, exists := s.IsExists(stud.Studio)
		if exists {
			continue
		}

		_, err := tx.Exec(`INSERT INTO studios (studio) SELECT ? WHERE NOT EXISTS(SELECT 1 FROM studios WHERE studio = ?)`, stud.Studio, stud.Studio)
		if err != nil {
			helpers.LogError(err.Error(), component)
			err := tx.Rollback()
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (s Model) Delete(studio string) {
	_, err := s.conn.Exec(`DELETE FROM tags WHERE tag = ?`, studio)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (s Model) Get(studiosQuery Studio) []Studio {
	query, args := models.QueryBuilderGet(studiosQuery, tableName)
	allStudios := make([]Studio, 0)

	row, err := s.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return allStudios
	}

	for row.Next() {
		studio := Studio{}
		err := row.Scan(&studio.GeneratedID, &studio.Studio, &studio.Count)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allStudios = append(allStudios, studio)
	}

	return allStudios
}
