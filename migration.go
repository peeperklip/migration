package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type migration struct {
	Sql     *sql.DB
	dialect string
}

func NewMigration(sql *sql.DB, dialect string) *migration {
	return &migration{Sql: sql, dialect: dialect}
}

//type DatabaseHandle interface {
// commented until I can find a way to make it possible \
// to inject instances from outside this package
//	createHandle() *sql.DB
//}

// func (mig migration) createHandle() *sql.DB {
func createHandle() *sql.DB {
	dbName := os.Getenv("DB_NAME")
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	//ssl mode should be verify-full
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", dbUsername, dbPassword, dbName)
	//_, err := sql.Open("postgres", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (mig migration) getRanMigrations() []int64 {

	var ranMigrations []int64

	result, err := mig.Sql.Query(QueryForRanMigrations(mig.dialect))

	if err != nil {
		return ranMigrations
	}

	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	for result.Next() {
		var currentMigration int64
		_ = result.Scan(&currentMigration)
		{
			ranMigrations = append(ranMigrations, currentMigration)
		}
	}

	return ranMigrations

}
func (mig migration) prepareMigrationsTable() {
	_, _ = mig.Sql.Exec(GetCreateTableByDialect(mig.dialect))
}
func (mig migration) HasMigrationRan(migrationToCheck string) bool {
	for _, item := range mig.getRanMigrations() {
		if strconv.FormatInt(item, 10) == migrationToCheck {
			return true
		}
	}

	return false

}
func (mig migration) GenerateMigration() {
	currentTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	migrationDirName := fmt.Sprintf("migrations/%s", currentTimestamp)

	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		_ = os.Mkdir("migrations", 0771)
	}
	if _, err := os.Stat(migrationDirName); os.IsNotExist(err) {
		_ = os.Mkdir(migrationDirName, 0771)
	}

	_, _ = os.Create(fmt.Sprintf("%s/up.sql", migrationDirName))
	_, _ = os.Create(fmt.Sprintf("%s/down.sql", migrationDirName))
}
func (mig migration) RunMigrations() {
	dir, err := os.ReadDir("migrations")
	if err != nil {
		return
	}

	dbUtil := DBUtil{}
	dbUtil.bootstrap()

	regexFofMigrationFile := regexp.MustCompile("\\d+")

	//breaks if the first migrations drops and recreates the database
	for _, value := range dir {

		if !value.IsDir() {
			continue
		}

		if !regexFofMigrationFile.MatchString(value.Name()) {
			continue
		}

		if mig.HasMigrationRan(value.Name()) == false {

			migrationFile, _ := os.ReadFile(fmt.Sprintf("%s/%s/up.sql", "migrations", value.Name()))

			fmt.Println("Currently executing: " + value.Name())
			_, err := mig.Sql.Exec(string(migrationFile))
			if err != nil {
				fmt.Println("Error when running migration: ")
				fmt.Println(err)
			}

			_, err = mig.Sql.Exec(InsertNewEntry(mig.dialect), value.Name())

			if err != nil {
				fmt.Println("when marking the migration as ran: ")
				fmt.Println(err)
			}
		}
	}
}

func (mig migration) closeConnection() {

}
