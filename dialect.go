package main

func GetCreateTableByDialect(dialect string) string {
	switch dialect {
	case "postgress":
		return "CREATE TABLE IF NOT EXISTS main.migrations (migration INTEGER NOT NULL);"
	case "test":
		return "NONE FOR NOW"
	default:
		panic("Could not figure out how to set up the migrations table")
	}
	return ""
}

func InsertNewEntry(dialect string) string {
	switch dialect {
	case "postgress":
		return "INSERT INTO main.migrations VALUES ($1);"

	case "test":
		return "NONE FOR NOW"
	default:
		panic("Could not figure out how to mark this migration as ran")
	}
	return ""
}

func QueryForRanMigrations(dialect string) string {
	switch dialect {
	case "postgress":
		return "SELECT migration FROM main.migrations"

	case "test":
		return "NONE FOR NOW"
	default:
		panic("Could not figure out how to mark this migration as ran")
	}
	return ""
}
