package auth

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry_go/pkg/database"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	tableName = "auth"
	component = "authModel"
)

type Auth struct {
	Username string `row:"username" type:"exact"`
	Password string `row:"password" type:"exact"`
	IsAdmin  bool   `row:"isadmin" type:"exact"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (a Model) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (a Model) Create(auth Auth) {
	auth.Password = hashPassword(auth.Password)

	query, args := models.QueryBuilderCreate(auth, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (a Model) Get(auth Auth) []Auth {
	query, args := models.QueryBuilderGet(auth, tableName)
	result := make([]Auth, 0)

	rows, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return result
	}

	for rows.Next() {
		scan := Auth{}
		err := rows.Scan(&scan.Username, &scan.Password, &scan.IsAdmin)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		result = append(result, scan)
	}
	return result
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	return string(hash)
}
