package auth

import (
	"database/sql"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	tableName = "auth"
	component = "authModel"
)

type Auth struct {
	Username string `row:"username" type:"exact"`
	Password string `row:"password" type:"exact"`
}

type AuthModel struct {
	conn *sql.DB
}

func Initialize() *AuthModel {
	return &AuthModel{
		conn: database.GetConn(router.GetDatabase(tableName)),
	}
}

func (a AuthModel) Create(auth Auth) {
	auth.Password = hashPassword(auth.Password)

	query, args := models.QueryBuilderCreate(auth, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}

func (a AuthModel) Get(auth Auth) []Auth {
	query, args := models.QueryBuilderGet(auth, tableName)
	result := make([]Auth, 0)

	rows, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return result
	}

	for rows.Next() {
		scan := Auth{}
		err := rows.Scan(&scan.Username, &scan.Password)
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