package auth

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	tableName = "auth"
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
		conn: models.GetConn(tableName),
	}
}

func (a Model) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (a Model) Create(auth Auth) {
	auth.Password = hashPassword(auth.Password)
	query, args := models.QueryBuilderCreate(auth, tableName)

	fmt.Println(args)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (a Model) Get(auth Auth) (allAuth []Auth) {
	query, args := models.QueryBuilderGet(auth, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &allAuth)
	return
}

// Hashes password string for storage
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Cost: 10
	if err != nil {
		helpers.LogError(err.Error())
	}

	return string(hash)
}
