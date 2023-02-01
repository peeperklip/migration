package migrations

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
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

	mig.RunMigrations()

	err = os.RemoveAll("migrations")
	err = os.RemoveAll("database.db")

	if err != nil {
		t.Error("Could not clean up after")
	}
}
