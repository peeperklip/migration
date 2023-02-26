package migrations

func GetCreateTableByDialect(dialect string) string {
	switch dialect {
	case "postgress":
		return "CREATE TABLE IF NOT EXISTS main.migrations (migConfig INTEGER NOT NULL);"
	case "sqlite3":
		return "CREATE TABLE IF NOT EXISTS migrations (migConfig INTEGER NOT NULL);"
	default:
		panic("Could not figure out how to set up the migrations table")
	}
	return ""
}

func InsertNewEntry(dialect string) string {
	switch dialect {
	case "postgress":
		return "INSERT INTO main.migrations VALUES ($1);"

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
