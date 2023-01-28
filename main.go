package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func Init(db *sql.DB) {
	args := os.Args
	if len(args) == 1 {
		contents, _ := os.ReadFile("help_file.txt")
		fmt.Print(string(contents))

		return
	}

	command := args[1]

	newMig := NewMigration(db, "postgress")

	switch command {

	case "generate":

		newMig.GenerateMigration()
		break

	case "migrate":
		newMig.RunMigrations()
		break

	case "down":

		break

	case "revert_to":

		break

	case "status":

		break

	case "":
		contents, _ := os.ReadFile("cli/help_file.txt")

		fmt.Println("This command is not supported \n" + string(contents))
	default:
		contents, _ := os.ReadFile("cli/help_file.txt")

		panic("This command is not supported \n" + string(contents))
	}
}
