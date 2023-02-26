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

	mig := migConfig{
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

	defer tearDown()
}

func TestMigration_GenerateMigration(t *testing.T) {
	mig, _ := setUp("testing_data/")

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

	defer tearDown()

}

func TestMigration_GetAllMigrations(t *testing.T) {
	mig, _ := setUp("")

	mig.GenerateMigration()
	mig.GenerateMigration()

	_ = exec.Command("cp", "--recursive", "testing_data", ".")

	migs := mig.GetAllMigrations()
	if len(migs) == 0 {
		t.Error("Expected at least one migConfig to be present on filesystem")
	}

	defer tearDown()
}

func TestMigration_GetUnRanMigrations(t *testing.T) {

	mig, _ := setUp("")
	mig.GenerateMigration()
	res := mig.GetUnRanMigrations()
	if len(res) != 1 {
		t.Error("failure!")
	}

	mig.RunMigrations()
	res = mig.GetUnRanMigrations()
	if len(res) != 0 {
		t.Error("failure!")
	}

	defer tearDown()
}

func tearDown() {
	_ = os.RemoveAll("migrations")
	_ = os.RemoveAll("database.db")
}

func setUp(baseDir string) (migConfig, error) {
	_, err := os.Create("database.db")
	db, err := sql.Open("sqlite3", "database.db")

	mig := migConfig{
		baseDir: baseDir,
		dialect: "sqlite3",
		Sql:     db,
	}

	return mig, err
}
