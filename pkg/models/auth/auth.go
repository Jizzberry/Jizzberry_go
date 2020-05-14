package auth

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry-go/pkg/models"
	"golang.org/x/crypto/bcrypt"
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
		conn: database.GetConn(router.GetDatabase("auth")),
	}
}

func (a AuthModel) Create(auth Auth) {
	auth.Password = hashPassword(auth.Password)

	query, args := models.QueryBuilderCreate(auth, "auth")

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
	}
}

func (a AuthModel) Get(auth Auth) []Auth {
	query, args := models.QueryBuilderGet(auth, "auth")
	result := make([]Auth, 0)

	rows, err := a.conn.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return result
	}

	for rows.Next() {
		scan := Auth{}
		err := rows.Scan(&scan.Username, &scan.Password)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, scan)
	}
	return result
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}
