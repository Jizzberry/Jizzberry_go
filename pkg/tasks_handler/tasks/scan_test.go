package tasks_test

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/database"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler/tasks"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestScan_Stop(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")

	runTestCancel(t, dir)

	//TODO: prevents runtime errors after temp dir is deleted (need better solution)
	time.Sleep(1 * time.Second)

	err := cleanDir(dir)
	if err != nil {
		t.Error(err)
	}
}

func TestScan_Start(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	runTestRW(t, dir)

	err := cleanDir(dir)
	if err != nil {
		t.Error(err)
	}
}

func runTestRW(t *testing.T, dir string) {
	helpers.SetWorkingDirectory(dir)
	helpers.CreateDirs()

	database.RunMigrations()

	err := createTmpFiles(dir, 10)
	if err != nil {
		t.Error(err)
	}

	tasks.Scan{}.Start([]string{dir})

	time.Sleep(5 * time.Second)

	err = verifyDatabase(10)
	if err != nil {
		t.Error(err)
	}
}

func runTestCancel(t *testing.T, dir string) {
	helpers.SetWorkingDirectory(dir)
	helpers.CreateDirs()

	database.RunMigrations()

	err := createTmpFiles(dir, 1)
	if err != nil {
		t.Error(err)
	}

	cancel, _ := tasks.Scan{}.Start([]string{dir})
	tmp := *cancel
	time.Sleep(3 * time.Second)
	tmp()
}

func cleanDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

func verifyDatabase(n int) error {
	conn := database.GetConn(router.GetDatabase("files"))
	defer conn.Close()
	rows, err := conn.Query(`SELECT count(generated_id) FROM files`)

	if err != nil {
		return err
	}
	defer rows.Close()
	rowsAffected := 0
	for rows.Next() {
		err := rows.Scan(&rowsAffected)
		if err != nil {
			return err
		}
	}

	if rowsAffected != n {
		return fmt.Errorf("not all files were scanned")
	}
	return nil
}

func createTmpFiles(dir string, n int) error {
	for i := 0; i < n; i++ {
		emptyFile, err := os.Create(filepath.FromSlash(dir + "/test" + strconv.Itoa(i) + ".mp4"))
		if err != nil {
			return err
		}
		emptyFile.Close()
	}
	return nil
}
