package migrations

type dbUtil struct {
	dbms            string
	createStatement string
}

var mapangpang = make(map[string]string)

func addCreateStatement(dbms string, createStatement string) {
	mapangpang[dbms] = createStatement
}

func GetCreateTableByDialect(dialect string) string {
	addCreateStatement("postgress", "CREATE TYPE migstatus AS ENUM ('RAN', 'REVERTED', 'UNRAN');CREATE TABLE IF NOT EXISTS main.migrations (migration INTEGER NOT NULL, migstatus migstatus);")

	return mapangpang[dialect]
}

func InsertNewEntry(dialect string) string {
	switch dialect {
	case "postgress":
		return "INSERT INTO main.migrations VALUES ($1, 'RAN');"

	case "sqlite3":
		return "INSERT INTO migrations VALUES ($1);"
	default:
		panic("Could not figure out how to mark this migConfig as ran")
	}
	return ""
}

func QueryForRanMigrations(dialect string) string {
	switch dialect {
	case "postgress":
		return "SELECT migConfig FROM main.migrations"

	case "sqlite3":
		return "SELECT migConfig FROM migrations"
	default:
		panic("Could not query for ran migrations")
	}
	return ""
}
