package studios

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry_go/pkg/database"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
)

const (
	tableName = "studios"
	component = "studiosModel"
)

type Studio struct {
	GeneratedID int64  `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	Name        string `row:"studio" type:"like" json:"generated_id"`
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

func (m Model) Close() {
	err := m.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (m Model) isEmpty() bool {
	rows, err := m.conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

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

func (m Model) IsExists(studio string) (int64, bool) {
	if m.isEmpty() {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return -1, false
	}

	fetch, err := m.conn.Query(`SELECT generated_id FROM studios WHERE studio = ?`, studio)
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

func (m Model) Create(studios []Studio) {
	tx, err := m.conn.Begin()

	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	for _, stud := range studios {
		_, exists := m.IsExists(stud.Name)
		if exists {
			continue
		}

		_, err := tx.Exec(`INSERT INTO studios (studio, count) SELECT ?, ? WHERE NOT EXISTS(SELECT 1 FROM studios WHERE studio = ?)`, stud.Name, stud.Count, stud.Name)
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

func (m Model) Delete(studio string) {
	_, err := m.conn.Exec(`DELETE FROM tags WHERE tag = ?`, studio)

	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (m Model) Get(studiosQuery Studio) []Studio {
	query, args := models.QueryBuilderGet(studiosQuery, tableName)
	allStudios := make([]Studio, 0)

	row, err := m.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return allStudios
	}

	for row.Next() {
		studio := Studio{}
		err := row.Scan(&studio.GeneratedID, &studio.Name, &studio.Count)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		allStudios = append(allStudios, studio)
	}

	return allStudios
}

func (m Model) GetFromTitle(names []string) []Studio {
	fetched := make([]Studio, 0)
	for _, name := range names {
		rows, err := m.conn.Query(`SELECT generated_id, studio FROM studios WHERE (studio LIKE ? COLLATE NOCASE) 
                                                         OR (replace(studio, ' ', '') LIKE ? COLLATE NOCASE)`, "%"+name+"%", name)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return fetched
		}

		for rows.Next() {
			var actor = Studio{}
			err := rows.Scan(&actor.GeneratedID, &actor.Name)
			if err != nil {
				helpers.LogError(err.Error(), component)
			}

			if !containsStudio(fetched, actor) {
				fetched = append(fetched, actor)
			}
		}
	}
	return fetched
}

func containsStudio(s []Studio, e Studio) bool {
	for _, a := range s {
		if a.GeneratedID == e.GeneratedID {
			return true
		}
	}
	return false
}
