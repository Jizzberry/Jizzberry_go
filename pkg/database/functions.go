package database

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/markbates/pkger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

func GetConn(databasePath string) *sql.DB {

	conn, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		helpers.LogError(err.Error())
	}

	return conn
}

func RunMigrations() error {
	dataDatabasepath := router.GetDatabase("files")
	actorsDatabasepath := router.GetDatabase("actors")
	studiosDatabasepath := router.GetDatabase("studios")
	authDatabasepath := router.GetDatabase("auth")

	migrationsData := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/jizzberry_data"),
	}

	migrationsActors := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/actors"),
	}

	migrationsStudios := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/studios"),
	}

	migrationsAuth := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/pkg/database/migrations/auth"),
	}

	err := doMigrate(migrationsData, dataDatabasepath)
	if err != nil {
		return err
	}
	err = doMigrate(migrationsActors, actorsDatabasepath)
	if err != nil {
		return err
	}
	err = doMigrate(migrationsAuth, authDatabasepath)
	if err != nil {
		return err
	}
	err = doMigrate(migrationsStudios, studiosDatabasepath)
	if err != nil {
		return err
	}
	return nil
}

func doMigrate(migrations *migrate.HttpFileSystemMigrationSource, databasePath string) error {
	conn := GetConn(databasePath)

	n, err := migrate.Exec(conn, "sqlite3", migrations, migrate.Up)

	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	helpers.LogInfo(fmt.Sprintf("Applied %d migrations in %s", n, databasePath))
	return nil
}
