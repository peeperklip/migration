package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func TestGenerateMigration(t *testing.T) {
	_, err := os.Create("database.db")
	if err != nil {
		t.Error("Could not create a database to start testing with")
	}
	db, err := sql.Open("sqlite3", "database.db")

	mig := migration{
		baseDir: "",
		dialect: "sqlite3",
		Sql:     db,
	}

	mig.GenerateMigration()

	validFilePath := regexp.MustCompile("^\\d*")

	dirs, _ := os.ReadDir("migrations")

	for _, val := range dirs {
		if val.IsDir() && val.Name() == "migrations" {
			continue
		}

		if !val.IsDir() {
			t.Error("Was not expecting a non directory here")
		}
		if !validFilePath.Match([]byte(val.Name())) {
			t.Error("Dir not match expected directory format")
		}

		sqlFiles, _ := os.ReadDir("migrations/" + val.Name())
		for _, sqlFile := range sqlFiles {
			if sqlFile.Name() != "down.sql" && sqlFile.Name() != "up.sql" {
				t.Error(fmt.Sprintf("%s expected, got %s", "up|down.sql", sqlFile.Name()))
			}
		}
	}

	mig.RunMigrations()

	defer func() {
		errMig := os.RemoveAll("migrations")
		errDb := os.RemoveAll("database.db")

		if errMig != nil || errDb != nil {
			t.Error("Could not clean up after")
		}
	}()
}

func TestMigration_GenerateMigration(t *testing.T) {
	_, err := os.Create("database.db")
	if err != nil {
		t.Error("Could not create a database to start testing with")
	}
	db, err := sql.Open("sqlite3", "database.db")

	mig := migration{
		baseDir: "testing_data/",
		dialect: "sqlite3",
		Sql:     db,
	}
	_ = exec.Command("cp", "--recursive", "testing_data", ".")

	mig.RunMigrations()

	result, err := mig.Sql.Query(QueryForRanMigrations(mig.dialect))

	if err != nil {
		t.Error("Failure")
		t.Error(err)
	}

	defer func(result *sql.Rows) {
		err = result.Close()
		if err != nil {
			t.Error(err)
		}
	}(result)

	defer func() {
		errMig := os.RemoveAll("migrations")
		errDb := os.RemoveAll("database.db")

		if errMig != nil || errDb != nil {
			t.Error("Could not clean up after")
		}
	}()

}
