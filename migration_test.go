package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peeperklip/migration/internal"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func TestGenerateMigration(t *testing.T) {
	tearDown()
	_, err := os.Create("database.db")
	if err != nil {
		t.Error("Could not create a database to start testing with")
	}

	db, err := sql.Open("sqlite3", "database.db")

	if err != nil {
		t.Error("Could not open DB")
	}
	mig := migConfig{
		baseDir: "",
		dialect: "sqlite3",
		sql:     db,
	}
	mig.initialize()

	mig.generateMigration()

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

	mig.runMigrations()
	tearDown()
}

func TestMigration_GenerateMigration(t *testing.T) {
	mig, _ := setUpTesting("testing_data/")

	_ = exec.Command("cp", "--recursive", "testing_data", ".")

	mig.runMigrations()

	result, err := mig.sql.Query(getQueryForGettingMigrations(mig.dialect))

	if err != nil {
		t.Error("Failure")
		t.Error(err)
		t.Error(internal.GetErrors())
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
	mig, _ := setUpTesting("")

	mig.generateMigration()
	mig.generateMigration()

	_ = exec.Command("cp", "--recursive", "testing_data", ".")

	migs := mig.getAllMigrations()
	if len(migs) == 0 {
		t.Error("Expected at least one migConfig to be present on filesystem")
	}

	defer tearDown()
}

func TestMigration_GetUnRanMigrations(t *testing.T) {

	mig, _ := setUpTesting("")
	mig.generateMigration()
	res := mig.getUnRanMigrations()
	if len(res) != 1 {
		t.Error("failure! expected 1")
	}

	mig.runMigrations()
	res = mig.getUnRanMigrations()

	if len(res) != 0 {
		t.Error("failure! expected 0")
		t.Error(len(res))
	}

	defer tearDown()
}
func tearDown() {
	_ = os.RemoveAll("migrations")
	_ = os.RemoveAll("database.db")
	internal.FlushErros()
}

func setUpTesting(baseDir string) (migConfig, error) {
	_, err := os.Create("database.db")
	db, err := sql.Open("sqlite3", "database.db")

	mg := NewMigration(db, "sqlite3", baseDir)

	mg.initialize()
	return *mg, err
}
