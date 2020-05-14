package database

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/database/router"
	"github.com/markbates/pkger"
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
	authDatabasepath := router.GetDatabase("auth")

	migrationsData := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/jizzberry_data"),
	}

	migrationsActors := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/actors"),
	}

	migrationsAuth := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/auth"),
	}

	doMigrate(migrationsData, dataDatabasepath)
	doMigrate(migrationsActors, actorsDatabasepath)
	doMigrate(migrationsAuth, authDatabasepath)

}

func doMigrate(migrations *migrate.HttpFileSystemMigrationSource, databasePath string) {
	conn := GetConn(databasePath)

	n, err := migrate.Exec(conn, "sqlite3", migrations, migrate.Up)

	if err != nil {
		fmt.Print(err)
	}

	conn.Close()
	fmt.Printf("Applied %d migrations!\n", n)

}
