package database

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	sqlite3mig "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

func getConn(databasePath string) *sql.DB {
	conn, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		fmt.Print("getConn(): %q\n", err)
	}
	return conn
}

func RunMigrations() {
	databasePath := filepath.FromSlash("../assets/database/jizzberry_data.db")

	conn := getConn(databasePath)

	driver, err := sqlite3mig.WithInstance(conn, &sqlite3mig.Config{})

	if err != nil {
		fmt.Printf("RunMigrations(): %s\n", err)
	}

	_, err = migrate.NewWithDatabaseInstance(
		"file://"+filepath.FromSlash("./migrations"),
		"ql",
		driver,
	)

	if err != nil {
		fmt.Printf("RunMigrations(): %s\n", err)
	}
}

func main() {
	RunMigrations()
}
