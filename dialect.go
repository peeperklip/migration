package migrations

import "fmt"

type dbUtil struct {
	dbms                  string
	createStatement       string
	insertStatement       string
	ranMigrationStatement string
}

var creationStatements = make(map[string]string)
var insertStatements = make(map[string]string)
var queryForRanMigrations = make(map[string]string)

func addCreateStatement(dbms string, createStatement string) {
	creationStatements[dbms] = createStatement
}

func addInsertStatement(dbms string, statement string) {
	insertStatements[dbms] = statement
}

func addQueryStatements(dbms string, statement string) {
	queryForRanMigrations[dbms] = statement
}

func setUp(dbms string) *dbUtil {
	for index, _ := range creationStatements {
		if index == dbms {
			break
		}

		panic(fmt.Sprintf("DBMS %s was not configured", dbms))
	}

	return &dbUtil{
		dbms:                  dbms,
		createStatement:       getCreateTableByDialect(dbms),
		insertStatement:       getInsertNewEntryByDialect(dbms),
		ranMigrationStatement: getQueryForGettingMigrations(dbms),
	}
}

func getCreateTableByDialect(dialect string) string {
	addCreateStatement("postgress", "CREATE TYPE migstatus AS ENUM ('RAN', 'REVERTED', 'UNRAN');CREATE TABLE IF NOT EXISTS main.migrations (migration INTEGER NOT NULL, migstatus migstatus);")
	addCreateStatement("sqlite3", "CREATE TABLE IF NOT EXISTS migrations (migration INTEGER NOT NULL, migstatus TEXT CHECK(migstatus IN('RAN', 'REVERTED', 'UNRAN')) NOT NULL DEFAULT 'UNRAN')")

	return creationStatements[dialect]
}

func getInsertNewEntryByDialect(dbms string) string {
	addInsertStatement("postgress", "INSERT INTO main.migrations VALUES ($1, 'RAN');")
	addInsertStatement("sqlite3", "INSERT INTO migrations VALUES ($1, 'RAN');")

	return insertStatements[dbms]
}

func getQueryForGettingMigrations(dbms string) string {
	addQueryStatements("postgress", "SELECT migration, migstatus FROM main.migrations")
	addQueryStatements("sqlite3", "SELECT migration, migstatus FROM migrations")

	return queryForRanMigrations[dbms]
}
