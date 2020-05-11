package database

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/gobuffalo/packr/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

func GetConn(databasePath string) *sql.DB {

	conn, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		fmt.Print("getConn():", err)
	}

	return conn
}

func RunMigrations() {
	dataDatabasepath := router.GetDatabase("files")
	actorsDatabasepath := router.GetDatabase("actors")

	migrationsData := &migrate.PackrMigrationSource{
		Box: packr.New("migrationsData", "./migrations/jizzberry_data"),
	}

	migrationsActors := &migrate.PackrMigrationSource{
		Box: packr.New("migrationsActors", "./migrations/actors"),
	}

	doMigrate(migrationsData, dataDatabasepath)
	doMigrate(migrationsActors, actorsDatabasepath)

}

func doMigrate(migrations *migrate.PackrMigrationSource, databasePath string) {
	conn := GetConn(databasePath)

	n, err := migrate.Exec(conn, "sqlite3", migrations, migrate.Up)

	if err != nil {
		fmt.Print(err)
	}

	conn.Close()
	fmt.Printf("Applied %d migrations!\n", n)

}
