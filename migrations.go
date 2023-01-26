package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		contents, _ := os.ReadFile("cli/help_file.txt")
		fmt.Print(string(contents))

		return
	}

	command := args[1]

	if command == "generate" {
		fmt.Println("Generating a migration ...")
		generateMigration()
		return
	}

	if command == "migrate" {
		fmt.Println("Running migrations (if needed)")
		runMigrations()
		return
	}

	if command == "down" {
		fmt.Println("")
	}

	contents, _ := os.ReadFile("cli/help_file.txt")

	panic("This command is not supported \n" + string(contents))
}

func runMigrations() {
	dir, err := os.ReadDir("migrations")
	if err != nil {
		return
	}

	regexFofMigrationFile := regexp.MustCompile("\\d+")

	ranMigrations := getRanMigrations()

	//breaks if the first migrations drops and recreates the database
	for _, value := range dir {

		if !value.IsDir() {
			continue
		}

		if !regexFofMigrationFile.MatchString(value.Name()) {
			continue
		}

		if hasMigrationRan(value.Name(), ranMigrations) == false || len(ranMigrations) == 0 {

			migration, _ := os.ReadFile(fmt.Sprintf("%s/%s/up.sql", "migrations", value.Name()))

			fmt.Println("Currently executing: " + value.Name())
			_, err := utils.CreateConnection().Exec(string(migration))
			if err != nil {
				fmt.Println("Error when running migration: ")
				fmt.Println(err)
			}

			prepareMigrationsTable()

			_, err = utils.CreateConnection().Exec("INSERT INTO main.migrations VALUES ($1)", value.Name())

			if err != nil {
				fmt.Println("when marking the migration as ran: ")
				fmt.Println(err)
			}
		}
	}
}

func generateMigration() {
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

func getRanMigrations() []int64 {
	var ranMigrations []int64

	result, err := utils.CreateConnection().Query("SELECT migration FROM main.migrations")

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

func prepareMigrationsTable() {
	_, err := utils.CreateConnection().Exec("create table If not exists main.migrations (migration integer not null);")
	if err != nil {
		fmt.Println(err)
	}
}

func hasMigrationRan(migration string, ranMigrations []int64) bool {
	for _, item := range ranMigrations {
		if strconv.FormatInt(item, 10) == migration {
			return true
		}
	}

	return false
}
